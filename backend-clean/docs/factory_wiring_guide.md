# Presenter／UseCase／Repository をファクトリで束ねる配線ガイド

Clean Architecture の「内側は外側を知らない」を守りつつ、依存の組み合わせを driver に集約するための方針をまとめます。RightCode 記事の Driver 配線パターン（`inputFactory` / `outputFactory` / `repositoryFactory`）をベースにしています。

## 何をどこで new するか
- **Controller**: 入力を受け、必要なファクトリを呼び出して UseCase に渡すだけ（ロジックは持たない）。
- **UseCase**: `InputPort` / `OutputPort` / `Repository` などのインターフェースにのみ依存し、具象は知らない。
- **Presenter**: `OutputPort` 実装。ドメインモデルをレスポンス DTO に詰め替え、結果を内部に保持する（リクエストごとに新規生成）。
- **Repository**: `Repository` インターフェース実装（DB などの外部依存をここに閉じ込める）。
- **Driver (cmd/api/main.go など)**: 具体型の生成と組み立て（DI）を一手に引き受ける。

## 推奨ファクトリ構成
```
internal/port
  account_port.go        // AccountInputPort, AccountOutputPort, AccountRepository
  note_port.go           // ...
  template_port.go       // ...

internal/adapter/http/presenter
  account_presenter.go   // func NewAccountPresenter() port.AccountOutputPort
  note_presenter.go
  template_presenter.go

internal/adapter/gateway/db
  account_repository.go  // func NewAccountRepository(...) port.AccountRepository
  ...

cmd/api/main.go (driver)
  type (
    AccountOutputFactory    func() port.AccountOutputPort
    AccountInputFactory     func(repo port.AccountRepository) port.AccountInputPort
    AccountRepositoryFactory func() port.AccountRepository
  )

  outputFactory := func() port.AccountOutputPort { return presenter.NewAccountPresenter() }
  repoFactory   := func() port.AccountRepository { return gatewaydb.NewAccountRepository(pool) }
  inputFactory  := func(repo port.AccountRepository) port.AccountInputPort {
    return usecase.NewAccountInteractor(repo)
  }
  ctrl := controller.NewAccountController(inputFactory, outputFactory, repoFactory)
```

## リクエストの流れ（例: Account）
```
[Controller] 受信 → output := outputFactory()
            → repo := repoFactory()
            → input := inputFactory(repo)
            → input.CreateOrGet(ctx, inputDTO, output)
[UseCase]    ドメインロジック＆永続化 → output.PresentAccount(ctx, account)
[Presenter]  ドメイン → OpenAPI DTO に詰め替え、内部に保持
[Controller] Presenter からレスポンス DTO を取り出し、HTTP で返す
```

## なぜファクトリ経由にするか
- **依存の集中**: 具体型の組み合わせを driver に閉じ込め、内側の層はインターフェースだけを知る。
- **テスト容易性**: ファクトリを差し替えてモックを渡せる。UseCase 単体テストでは OutputPort/Repository をモック化するだけでよい。
- **安全性**: Presenter は状態を持つためリクエストごとに new する（ファクトリなら使い回さない）。

## レイヤー責務の再確認
- **Domain**: エンティティ/VO/ドメインサービス。FW/DB 型を持ち込まない。
- **UseCase**: 入力を受けてドメインロジックを組み立て、OutputPort に結果を渡す。外部 I/F 型は知らない。
- **Interface Adapter**: Presenter でレスポンス変換、Gateway で DB/外部API 変換、Controller は配線のみ。
- **Driver**: FW 初期化・依存注入（ファクトリ）・設定。

この方針で、各ドメイン（Account/Note/Template）の Controller に `inputFactory` と `outputFactory` を渡し、UseCase は OutputPort を呼ぶだけの形に統一してください。
