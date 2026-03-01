package usecase

import (
	"context"
	"time"
	"app/src/domain/entities"
	"app/src/domain/services" // auth_domain_service.go を利用するために追加
)

type RegisterSpotPostInput struct {
	Token     string // トークンを受け取るように変更
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
	presenter   RegisterSpotPostPresenter
	spotRepo    entities.SpotRepository
	postRepo    entities.PostRepository
	authService services.AuthDomainService // 認証サービスを追加
}

func NewRegisterSpotPostInteractor(
	p RegisterSpotPostPresenter,
	s entities.SpotRepository,
	r entities.PostRepository,
	a services.AuthDomainService, // コンストラクタに認証サービスを追加
) RegisterSpotPostUseCase {
	return &registerSpotPostInteractor{
		presenter:   p,
		spotRepo:    s,
		postRepo:    r,
		authService: a,
	}
}

func (i *registerSpotPostInteractor) Execute(ctx context.Context, input RegisterSpotPostInput) (*RegisterSpotPostOutput, error) {
	// 1. ユーザーの特定（認証）
	// input.Token から user オブジェクトを特定
	user, err := i.authService.VerifyToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

	// 2. 座標による同一店舗の特定
	existingSpot, err := i.spotRepo.FindByLocation(
		ctx,
		input.Latitude,
		input.Longitude,
	)
	if err != nil {
		return nil, err
	}

	var targetSpot *entities.Spot

	// --- 【理想の挙動への変更】 ---
	// 既存店舗がある場合は、エラーを返さず自動的にその店舗を選択する
	if existingSpot != nil {
		targetSpot = existingSpot
	} else {
		// 【同一地点なし】新規 Spot を生成・保存
		// user.ID.Value() で int 型を取り出す
		newSpot, err := entities.NewSpot(
			0,
			input.SpotName,
			input.Latitude,
			input.Longitude,
			user.ID.Value(),
		)
		if err != nil {
			return nil, err
		}

		targetSpot, err = i.spotRepo.Create(newSpot)
		if err != nil {
			return nil, err
		}
	}

	// 3. Post（投稿）の生成
	// user.ID.Value() と user.Username.String() を使用して型を合わせる
	post, err := entities.NewPost(
		0,
		user.ID.Value(),
		targetSpot.ID.Value(),
		user.Username.String(),
		input.ImageURL,
		input.Caption,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	// 4. Post の永続化
	createdPost, err := i.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	// 5. 出力整形
	return i.presenter.Output(targetSpot, createdPost), nil
}