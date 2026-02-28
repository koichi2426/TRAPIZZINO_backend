package usecase

import (
	"context"
	"errors"
	"time"
	"app/domain/entities"
)

type RegisterSpotPostInput struct {
	UserID    int
	Username  string
	SpotName  string
	Latitude  float64
	Longitude float64
	ImageURL  string
	Caption   string
	Overwrite bool
}

type RegisterSpotPostOutput struct {
	SpotID int
	PostID int
}

type RegisterSpotPostPresenter interface {
	Output(spot *entities.Spot, post *entities.Post) *RegisterSpotPostOutput
}

type RegisterSpotPostUseCase interface {
	Execute(ctx context.Context, input RegisterSpotPostInput) (*RegisterSpotPostOutput, error)
}

type registerSpotPostInteractor struct {
	presenter RegisterSpotPostPresenter
	spotRepo  entities.SpotRepository
	postRepo  entities.PostRepository
}

func NewRegisterSpotPostInteractor(
	p RegisterSpotPostPresenter,
	s entities.SpotRepository,
	r entities.PostRepository,
) RegisterSpotPostUseCase {
	return &registerSpotPostInteractor{
		presenter: p,
		spotRepo:  s,
		postRepo:  r,
	}
}

func (i *registerSpotPostInteractor) Execute(ctx context.Context, input RegisterSpotPostInput) (*RegisterSpotPostOutput, error) {
	// 1. Spot の生成
	spot, err := entities.NewSpot(
		0,
		input.SpotName,
		input.Latitude,
		input.Longitude,
		input.UserID,
	)
	if err != nil {
		return nil, err
	}

	// 2. 既存スポットの確認
	// リポジトリのインターフェースに合わせて引数を MeshID のみに修正
	existing, err := i.spotRepo.FindByMeshID(spot.MeshID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 && !input.Overwrite {
		return nil, errors.New("このメッシュには既に登録があります。上書きを選択してください。")
	}

	// 3. Spot の永続化
	createdSpot, err := i.spotRepo.Create(spot)
	if err != nil {
		return nil, err
	}

	// 4. Post の生成
	// 修正ポイント：新しく定義した NewPost(id, userID, spotID, ...) の引数順に適合
	post, err := entities.NewPost(
		0,                       // id (新規作成時は 0)
		input.UserID,            // userID
		createdSpot.ID.Value(),  // spotID (直前で保存した Spot の ID を渡す)
		input.Username,          // username
		input.ImageURL,          // imageURL
		input.Caption,           // caption
		time.Now(),              // postedAt
	)
	if err != nil {
		return nil, err
	}

	// 5. Post の永続化
	// エンティティ生成時に SpotID がセットされているため、そのまま保存
	createdPost, err := i.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	// 6. 出力整形
	return i.presenter.Output(createdSpot, createdPost), nil
}