package usecase

import (
    "context"
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
    // 1. 座標による同一店舗の特定
    // 店名に関わらず、同じ座標に既に登録があるか確認します
    existingSpot, err := i.spotRepo.FindByLocation(
        ctx, 
        input.Latitude, 
        input.Longitude,
    )
    if err != nil {
        return nil, err
    }

    var targetSpot *entities.Spot
    if existingSpot != nil {
        // 【同一地点あり】既存の SpotID を使用（店舗の重複を防ぐ）
        targetSpot = existingSpot
    } else {
        // 【同一地点なし】新規 Spot を生成・保存
        newSpot, err := entities.NewSpot(
            0,
            input.SpotName,
            input.Latitude,
            input.Longitude,
            input.UserID,
        )
        if err != nil {
            return nil, err
        }
        
        targetSpot, err = i.spotRepo.Create(newSpot)
        if err != nil {
            return nil, err
        }
    }

    // 2. Post（投稿）の生成
    // 既存・新規どちらの場合も、確定した targetSpot.ID に紐付けます
    post, err := entities.NewPost(
        0,
        input.UserID,
        targetSpot.ID.Value(), 
        input.Username,
        input.ImageURL,
        input.Caption,
        time.Now(),
    )
    if err != nil {
        return nil, err
    }

    // 3. Post の永続化
    createdPost, err := i.postRepo.Create(post)
    if err != nil {
        return nil, err
    }

    // 4. 出力整形
    return i.presenter.Output(targetSpot, createdPost), nil
}