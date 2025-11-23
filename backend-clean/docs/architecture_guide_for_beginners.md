# 🏛️ クリーンアーキテクチャ入門ガイド - 本当にゼロから理解する

このドキュメントは、**小学生でもわかる**をコンセプトに、クリーンアーキテクチャを実例とともに解説します。

「難しそう…」と思った人、安心してください。このガイドを読めば、**なぜこう設計するのか**、**どこに何を書けばいいのか**が手に取るようにわかります。

---

## 📚 目次

1. [クリーンアーキテクチャって何？](#クリーンアーキテクチャって何)
2. [最も大切な3つのルール](#最も大切な3つのルール)
3. [家づくりで理解するレイヤー構造](#家づくりで理解するレイヤー構造)
4. [リクエストの冒険 - データの流れを追う](#リクエストの冒険---データの流れを追う)
5. [各レイヤーの詳細解説](#各レイヤーの詳細解説)
6. [実例で学ぶ: ノート作成機能](#実例で学ぶ-ノート作成機能)
7. [よくある質問と答え](#よくある質問と答え)
8. [チェックリスト: コードを書く前に](#チェックリスト-コードを書く前に)

---

## クリーンアーキテクチャって何？

クリーンアーキテクチャは、**「変わるもの」と「変わらないもの」を分ける設計手法**です。

### なぜ必要なの？

想像してみてください。あなたがレストランのオーナーだとします。

```
┌─────────────────────────────────────────┐
│  🍕 ピザレストラン                        │
│                                         │
│  ❌ 悪い設計：                            │
│  - キッチンで注文も受ける                 │
│  - レジで料理も作る                       │
│  - ウェイターがレシピを決める             │
│  → 混沌！誰が何をやるかわからない！        │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  🍕 ピザレストラン                        │
│                                         │
│  ✅ 良い設計：                            │
│  - ホール: お客さんから注文を受ける        │
│  - キッチン: レシピに従って料理を作る      │
│  - レジ: お会計を処理する                 │
│  → 役割が明確！変更も簡単！               │
└─────────────────────────────────────────┘
```

クリーンアーキテクチャは、**システムを「役割ごとの部屋」に分ける**設計です。

---

## 最も大切な3つのルール

このプロジェクトでは、3つのシンプルなルールを守るだけでOKです。

### 📌 ルール1: 依存は内向きだけ

```
外側のレイヤーは内側を知ってOK
内側のレイヤーは外側を知っちゃダメ

      ┌─────────────┐
      │   外側      │ ──┐
      │ (Controller)│   │ OK: 内側を知ってる
      └─────────────┘   ↓
            │
            ↓
      ┌─────────────┐
      │   内側      │   ×  NG: 外側を知らない
      │  (UseCase)  │      (外に矢印が出ない)
      └─────────────┘
```

**具体例：**
- ✅ Controller が UseCase を呼ぶ → OK
- ❌ UseCase が Controller を呼ぶ → NG
- ✅ Gateway が Domain の型を使う → OK
- ❌ Domain が PostgreSQL の型を使う → NG

### 📌 ルール2: データ構造を跨ぐときは変換する

```
HTTPリクエスト → ドメインモデル → DBの行

それぞれ別の型！変換が必要！

┌──────────────┐  変換   ┌──────────────┐  変換   ┌──────────────┐
│ HTTP JSON    │ ─────→ │ Domain Note  │ ─────→ │ PostgreSQL   │
│ (OpenAPI型)  │ Controller│ (純粋なGo)  │ Gateway│ (pgtype.UUID)│
└──────────────┘        └──────────────┘        └──────────────┘
```

**なぜ？**
- ドメインは**ビジネスルール**に集中したい
- フレームワークやDBが変わっても**ドメインは変えたくない**

### 📌 ルール3: UseCase は1メソッド1仕事

```
❌ 悪い例：
func (u *UseCase) DoEverything() {
  // ユーザー作成
  // 通知送信
  // レポート生成
  // ...何でもやる
}

✅ 良い例：
func (u *UseCase) CreateNote() { /* ノート作成だけ */ }
func (u *UseCase) PublishNote() { /* 公開だけ */ }
func (u *UseCase) DeleteNote() { /* 削除だけ */ }
```

**なぜ？**
- テストしやすい
- 変更の影響範囲が小さい
- トランザクション境界が明確

---

## 家づくりで理解するレイヤー構造

クリーンアーキテクチャを**家づくり**で例えます。

```
                    🏠 クリーンアーキテクチャの家

┌─────────────────────────────────────────────────────────┐
│                     🚪 玄関 (Cmd)                        │ ← アプリの入口
│                  ・訪問者を迎える                         │
│                  ・家の準備をする                         │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│               🔧 電気・水道 (Driver)                      │ ← インフラの配線
│          ・DB接続、設定、Factoryなど                      │
│          ・各部屋に必要なものを供給                       │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                 📞 受付窓口 (Controller)                  │ ← 外部とのやり取り
│          ・お客さんの要望を聞く                           │
│          ・日本語を家の人に通訳                           │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                🎯 指令室 (UseCase)                       │ ← 仕事の手順
│          ・何をどう処理するか決める                        │
│          ・各部屋に指示を出す                             │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│              ❤️ 心臓部 (Domain)                          │ ← ビジネスルール
│          ・家の根本的なルール                             │
│          ・「何が正しいか」を知ってる                      │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│         📦 倉庫・郵便局 (Gateway/Presenter)                │ ← 外部との変換
│     ・家のものを外の形式に変換                            │
│     ・外のものを家の形式に変換                            │
└─────────────────────────────────────────────────────────┘
```

### 各部屋の役割

| 部屋 | 実際のレイヤー | やること |
|------|--------------|----------|
| 🚪 玄関 | **Cmd** | main関数。サーバー起動だけ |
| 🔧 電気・水道 | **Driver** | DB接続、設定、Factory（配線） |
| 📞 受付窓口 | **Controller** | HTTPリクエストを受け取る |
| 🎯 指令室 | **UseCase** | ビジネスロジックの手順 |
| ❤️ 心臓部 | **Domain** | ビジネスルール（不変条件） |
| 📦 倉庫・郵便局 | **Gateway/Presenter** | DB/HTTPとの変換 |
| 📝 契約書 | **Port** | 各レイヤー間のインターフェース |

---

## リクエストの冒険 - データの流れを追う

実際に **GET /api/notes/123** というリクエストが来たら、何が起こるのか、冒険形式で追いかけましょう！

```
🌍 インターネット
    │
    │ GET /api/notes/123
    ↓
┌────────────────────────────────────────────────┐
│  1️⃣ Echo (Webフレームワーク)                   │
│     「おっ、リクエストが来た！」                │
└────────────────────────────────────────────────┘
    │
    │ 「/api/notes/:id のControllerを呼ぶよ」
    ↓
┌────────────────────────────────────────────────┐
│  2️⃣ Controller (受付窓口)                      │
│     func (c *NoteController) GetByID(...)      │
│                                                │
│     やること：                                  │
│     - パスパラメータ "123" を取り出す           │
│     - UseCaseに「ID=123のノート欲しい」と伝える │
└────────────────────────────────────────────────┘
    │
    │ input.Get(ctx, "123")
    ↓
┌────────────────────────────────────────────────┐
│  3️⃣ UseCase (指令室)                           │
│     func (u *NoteInteractor) Get(...)          │
│                                                │
│     やること：                                  │
│     - Repository に「ノート取ってきて」         │
│     - 取得したデータをPresenterに渡す           │
└────────────────────────────────────────────────┘
    │                              ↑
    │ Get(ctx, "123")              │
    ↓                              │ 取得したNote
┌──────────────────────┐           │
│  4️⃣ Repository       │           │
│  (Gateway - 倉庫)    │───────────┘
│                      │
│  やること：           │
│  - DBに問い合わせ     │
│  - DB行をDomainに変換 │
└──────────────────────┘
    │
    │ SELECT * FROM notes WHERE id = '123'
    ↓
┌────────────────────────────────────────────────┐
│  5️⃣ PostgreSQL                                 │
│     データベース                                │
│                                                │
│     返却: Row{id, title, ...}                  │
└────────────────────────────────────────────────┘
    │
    │ pgx.Row → note.WithMeta (ドメインモデル)
    ↓
┌────────────────────────────────────────────────┐
│  6️⃣ Presenter (郵便局)                         │
│     func (p *NotePresenter) PresentNote(...)   │
│                                                │
│     やること：                                  │
│     - note.WithMeta → OpenAPI Response         │
│     - CreatedAt を ISO8601形式にフォーマット    │
└────────────────────────────────────────────────┘
    │
    │ OpenAPI Response JSON
    ↓
┌────────────────────────────────────────────────┐
│  7️⃣ Controller が JSON返却                     │
│     return ctx.JSON(http.StatusOK, p.Note())   │
└────────────────────────────────────────────────┘
    │
    │ HTTP 200 OK
    ↓
🌍 インターネット (ユーザーのブラウザ)
```

### ポイント整理

```
🔄 データ変換の流れ

HTTP JSON ──Controller──→ string "123"
                ↓
         UseCase が処理
                ↓
         Repository が呼ばれる
                ↓
pgx.Row ──Gateway──→ note.WithMeta (Domain)
                ↓
         UseCase が Presenter に渡す
                ↓
note.WithMeta ──Presenter──→ OpenAPI Response
                ↓
         Controller が返す
                ↓
          HTTP JSON
```

**なぜこんなに変換するの？**

→ **各レイヤーが独立していて、フレームワークやDBを簡単に変更できるから！**

---

## 各レイヤーの詳細解説

### 1️⃣ Domain (ドメイン層) - ビジネスルールの中心

**場所:** `internal/domain/`

**役割:** アプリケーションの「心臓部」。ビジネスルールを表現する。

#### 何を書く？

```go
// ✅ Entity (エンティティ) - ビジネスの「もの」
type Note struct {
    ID         string
    Title      string
    TemplateID string
    OwnerID    string
    Status     NoteStatus  // Draft or Publish
    Sections   []Section
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// ✅ Value Object (値オブジェクト) - 状態
type NoteStatus string

const (
    StatusDraft   NoteStatus = "Draft"
    StatusPublish NoteStatus = "Publish"
)

// ✅ ビジネスルール - 「〜してはいけない」
func CanChangeStatus(from, to NoteStatus) error {
    if from == StatusDraft && to == StatusPublish {
        return nil  // OK: 下書き → 公開
    }
    if from == StatusPublish && to == StatusDraft {
        return nil  // OK: 公開 → 下書き
    }
    if from == to {
        return nil  // OK: 同じ状態
    }
    return domainerr.ErrInvalidStatusChange  // NG!
}

// ✅ 検証ロジック
func ValidateNoteOwnership(noteOwnerID, actorID string) error {
    if noteOwnerID != actorID {
        return domainerr.ErrUnauthorized
    }
    return nil
}
```

#### ❌ 書いちゃダメなもの

```go
// ❌ NG: HTTP関連
import "github.com/labstack/echo/v4"

// ❌ NG: DB関連
import "github.com/jackc/pgx/v5"

// ❌ NG: 外部フレームワーク
import openapi "..."
```

**ドメインは純粋なGoで書く！**

#### ファイル構成例

```
internal/domain/note/
├── entity.go           ← Note, Section の定義
├── types.go            ← NoteStatus などの値オブジェクト
├── logic.go            ← ビジネスルールの関数
├── aggregate.go        ← Noteと関連データをまとめた型
└── logic_test.go       ← テスト
```

---

### 2️⃣ UseCase (ユースケース層) - アプリケーションの手順書

**場所:** `internal/usecase/`

**役割:** 「何をどの順番でやるか」を決める指令室。

#### 何を書く？

```go
// ✅ UseCase (Interactor) の構造
type NoteInteractor struct {
    notes     port.NoteRepository      // ← Port経由で依存
    templates port.TemplateRepository
    tx        port.TxManager
    output    port.NoteOutputPort
}

// ✅ ビジネスロジックの手順
func (u *NoteInteractor) Create(ctx context.Context, input port.NoteCreateInput) error {
    // 1. テンプレート取得
    tpl, err := u.templates.Get(ctx, input.TemplateID)
    if err != nil {
        return err
    }

    // 2. ドメインルールで検証
    sections := buildSections("", tpl.Template.Fields, input.Sections)
    if err := note.ValidateNoteForCreate(input.Title, tpl.Template, sections); err != nil {
        return err
    }

    // 3. トランザクション開始
    var noteID string
    err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        // 4. ドメインサービスでノート構築
        built, err := service.BuildNote(input.Title, tpl.Template, input.OwnerID, sections)
        if err != nil {
            return err
        }

        // 5. Repository で保存
        nn, err := u.notes.Create(txCtx, built)
        if err != nil {
            return err
        }
        noteID = nn.ID

        // 6. セクション保存
        sectionsWithID := buildSections(noteID, tpl.Template.Fields, input.Sections)
        return u.notes.ReplaceSections(txCtx, noteID, sectionsWithID)
    })

    if err != nil {
        return err
    }

    // 7. 最新データ取得してPresenterに渡す
    n, err := u.notes.Get(ctx, noteID)
    if err != nil {
        return err
    }
    return u.output.PresentNote(ctx, n)
}
```

#### ポイント

- **手順書みたいに読める** ← コメントで番号振ってOK
- **Port（インターフェース）経由で外部を呼ぶ** ← 依存を逆転
- **トランザクション境界を決める** ← ここで開始・終了

#### ❌ 書いちゃダメなもの

```go
// ❌ NG: HTTPレスポンスを直接返す
return ctx.JSON(200, note)

// ❌ NG: 具体的なDB実装に依存
notes := &db.NoteRepository{}

// ✅ OK: Portに依存
notes port.NoteRepository
```

---

### 3️⃣ Port (ポート層) - 契約書

**場所:** `internal/port/`

**役割:** レイヤー間の「約束事」を定義するインターフェース。

#### 何を書く？

```go
// ✅ InputPort - UseCaseへの入口
type NoteInputPort interface {
    List(ctx context.Context, filters note.Filters) error
    Get(ctx context.Context, id string) error
    Create(ctx context.Context, input NoteCreateInput) error
    Update(ctx context.Context, input NoteUpdateInput) error
    Delete(ctx context.Context, id, ownerID string) error
}

// ✅ OutputPort - UseCaseからの出口（Presenterへ）
type NoteOutputPort interface {
    PresentNoteList(ctx context.Context, notes []note.WithMeta) error
    PresentNote(ctx context.Context, note *note.WithMeta) error
    PresentNoteDeleted(ctx context.Context) error
}

// ✅ Repository - データ永続化の契約
type NoteRepository interface {
    List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error)
    Get(ctx context.Context, id string) (*note.WithMeta, error)
    Create(ctx context.Context, n note.Note) (*note.Note, error)
    Update(ctx context.Context, n note.Note) (*note.Note, error)
    Delete(ctx context.Context, id string) error
}

// ✅ 入力DTO
type NoteCreateInput struct {
    Title      string
    TemplateID string
    OwnerID    string
    Sections   []SectionInput
}
```

#### なぜ必要？

```
┌──────────────┐                    ┌──────────────┐
│  UseCase     │                    │  UseCase     │
│              │                    │              │
│  具体的な実装│ ← これだとテスト   │  Interface   │ ← テスト簡単！
│  に依存      │    しにくい        │  に依存      │    Mockできる
└──────────────┘                    └──────────────┘
        ↓                                   ↓
┌──────────────┐                    ┌──────────────┐
│PostgresRepo  │                    │   Port       │
└──────────────┘                    └──────────────┘
                                            ↑
                                    ┌───────┴────────┐
                                    │                │
                              PostgresRepo      MockRepo
                              (本番)            (テスト)
```

---

### 4️⃣ Controller (コントローラー層) - HTTP入力の受付

**場所:** `internal/adapter/http/controller/`

**役割:** HTTPリクエストを受け取り、UseCaseに渡す「通訳」。

#### 何を書く？

```go
// ✅ Controller構造
type NoteController struct {
    inputFactory    func(...) port.NoteInputPort
    outputFactory   func() *presenter.NotePresenter
    noteRepoFactory func() port.NoteRepository
    tplRepoFactory  func() port.TemplateRepository
    txFactory       func() port.TxManager
}

// ✅ HTTPハンドラ
func (c *NoteController) Create(ctx echo.Context) error {
    // 1. リクエストボディをパース
    var body openapi.ModelsCreateNoteRequest
    if err := ctx.Bind(&body); err != nil {
        return ctx.JSON(http.StatusBadRequest,
            openapi.ModelsBadRequestError{
                Code:    openapi.ModelsBadRequestErrorCodeBADREQUEST,
                Message: "invalid body",
            })
    }

    // 2. OpenAPI型 → ドメイン入力DTOに変換
    ownerID := body.OwnerId.String()
    sections := []port.SectionInput{}
    if body.Sections != nil {
        for _, s := range *body.Sections {
            sections = append(sections, port.SectionInput{
                FieldID: s.FieldId,
                Content: s.Content,
            })
        }
    }

    // 3. UseCaseを呼ぶ
    input, p := c.newIO()
    err := input.Create(ctx.Request().Context(), port.NoteCreateInput{
        Title:      body.Title,
        TemplateID: body.TemplateId.String(),
        OwnerID:    ownerID,
        Sections:   sections,
    })
    if err != nil {
        return handleError(ctx, err)
    }

    // 4. Presenterから結果を取得して返す
    return ctx.JSON(http.StatusOK, p.Note())
}
```

#### ポイント

- **OpenAPI型 → ドメインDTOへの変換** ← Controllerの仕事
- **エラーハンドリング** ← ドメインエラーをHTTPステータスに変換
- **Factory経由でUseCaseを取得** ← DI（依存性注入）

---

### 5️⃣ Presenter (プレゼンター層) - HTTP出力の整形

**場所:** `internal/adapter/http/presenter/`

**役割:** ドメインモデルをHTTPレスポンスに変換する「翻訳者」。

#### 何を書く？

```go
// ✅ Presenter構造
type NotePresenter struct {
    note      *openapi.ModelsNoteResponse
    notes     []openapi.ModelsNoteResponse
    deletedOK bool
}

// ✅ OutputPortの実装
func (p *NotePresenter) PresentNote(_ context.Context, n *note.WithMeta) error {
    resp := toNoteResponse(*n)
    p.note = &resp  // 結果を保存
    return nil
}

// ✅ ドメインモデル → OpenAPI型への変換
func toNoteResponse(n note.WithMeta) openapi.ModelsNoteResponse {
    sections := make([]openapi.ModelsSection, 0, len(n.Sections))
    for _, s := range n.Sections {
        sections = append(sections, openapi.ModelsSection{
            Id:         s.Section.ID,
            FieldId:    s.Section.FieldID,
            FieldLabel: s.FieldLabel,
            Content:    s.Section.Content,
            IsRequired: s.IsRequired,
        })
    }

    return openapi.ModelsNoteResponse{
        Id:           n.Note.ID,
        Title:        n.Note.Title,
        TemplateId:   n.Note.TemplateID,
        TemplateName: n.TemplateName,
        OwnerId:      n.Note.OwnerID,
        Owner: openapi.ModelsAccountSummary{
            Id:        n.Note.OwnerID,
            FirstName: n.OwnerFirstName,
            LastName:  n.OwnerLastName,
            Thumbnail: n.OwnerThumbnail,
        },
        Status:    openapi.ModelsNoteStatus(n.Note.Status),
        Sections:  sections,
        CreatedAt: n.Note.CreatedAt,
        UpdatedAt: n.Note.UpdatedAt,
    }
}
```

#### ポイント

- **結果を内部で保存** ← Controllerが後で取り出す
- **ドメイン → OpenAPIへの変換** ← Presenterの仕事

---

### 6️⃣ Gateway (ゲートウェイ層) - DBとのやり取り

**場所:** `internal/adapter/gateway/db/`

**役割:** データベースからデータを取得し、ドメインモデルに変換する「倉庫番」。

#### 何を書く？

```go
// ✅ Repository実装
type NoteRepository struct {
    pool    *pgxpool.Pool
    queries *sqldb.Queries  // sqlc生成コード
}

// ✅ Portの実装
func (r *NoteRepository) Get(ctx context.Context, id string) (*note.WithMeta, error) {
    // 1. UUIDに変換
    pgID, err := toUUID(id)
    if err != nil {
        return nil, err
    }

    // 2. sqlcでDB問い合わせ
    row, err := queriesForContext(ctx, r.queries).GetNoteByID(ctx, pgID)
    if err != nil {
        // ❌ pgx.ErrNoRows をそのまま返さない！
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainerr.ErrNotFound  // ✅ ドメインエラーに変換
        }
        return nil, err
    }

    // 3. Sectionを取得
    sections, err := r.listSections(ctx, row.ID)
    if err != nil {
        return nil, err
    }

    // 4. DB行 → ドメインモデルに変換
    return &note.WithMeta{
        Note: note.Note{
            ID:         uuidToString(row.ID),      // pgtype.UUID → string
            Title:      row.Title,
            TemplateID: uuidToString(row.TemplateID),
            OwnerID:    uuidToString(row.OwnerID),
            Status:     note.NoteStatus(row.Status),
            CreatedAt:  timestamptzToTime(row.CreatedAt),
            UpdatedAt:  timestamptzToTime(row.UpdatedAt),
        },
        TemplateName:   row.TemplateName,
        OwnerFirstName: row.FirstName,
        OwnerLastName:  row.LastName,
        OwnerThumbnail: thumbnail,
        Sections:       sections,
    }, nil
}
```

#### ポイント

- **DB型 → ドメイン型への変換** ← `pgtype.UUID` → `string`
- **DBエラー → ドメインエラーへの変換** ← `pgx.ErrNoRows` → `domainerr.ErrNotFound`
- **sqlc生成コードを使う** ← タイプセーフなSQL

---

### 7️⃣ Driver (ドライバー層) - 配線工場

**場所:** `internal/driver/`

**役割:** DB接続、設定、Factory関数など「インフラの準備」。

#### 何を書く？

```go
// ✅ Factory関数
// internal/driver/factory/usecase_factory.go
func NewNoteInputFactory() func(
    noteRepo port.NoteRepository,
    tplRepo port.TemplateRepository,
    tx port.TxManager,
    output port.NoteOutputPort,
) port.NoteInputPort {
    return func(noteRepo port.NoteRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteOutputPort) port.NoteInputPort {
        return usecase.NewNoteInteractor(noteRepo, tplRepo, tx, output)
    }
}

// ✅ DB接続
// internal/driver/db/pool.go
func NewPool(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, error) {
    poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }
    return pgxpool.NewWithConfig(ctx, poolCfg)
}

// ✅ Initializer - 全体の組み立て
// internal/driver/initializer/api/initializer.go
func BuildServer(ctx context.Context) (*echo.Echo, func(), error) {
    // 1. 設定読み込み
    cfg := config.Load()

    // 2. DB接続
    pool, err := db.NewPool(ctx, cfg.DB)
    if err != nil {
        return nil, nil, err
    }

    // 3. Factoryを作成
    noteInputFactory := factory.NewNoteInputFactory()
    noteOutputFactory := factory.NewNotePresenterFactory()
    // ...

    // 4. Controllerを作成
    noteCtrl := controller.NewNoteController(
        noteInputFactory,
        noteOutputFactory,
        noteRepoFactory,
        tplRepoFactory,
        txFactory,
    )

    // 5. ルーティング設定
    server := controller.NewServer(noteCtrl, tplCtrl, accountCtrl)

    cleanup := func() {
        pool.Close()
    }

    return server.Echo, cleanup, nil
}
```

---

### 8️⃣ Cmd (コマンド層) - エントリーポイント

**場所:** `cmd/api/main.go`

**役割:** アプリケーションの起動。これだけ！

#### 何を書く？

```go
func main() {
    ctx := context.Background()

    // 1. サーバー構築
    e, cleanup, err := initializer.BuildServer(ctx)
    if err != nil {
        log.Fatalf("failed to initialize server: %v", err)
    }
    defer cleanup()

    // 2. サーバー起動
    addr := ":8080"
    log.Printf("starting HTTP server at %s\n", addr)
    if err := e.Start(addr); err != nil {
        log.Fatalf("server exited: %v", err)
    }
}
```

**main関数は短く！** ← 詳細はInitializerに任せる

---

## 実例で学ぶ: ノート作成機能

実際の「ノート作成」機能を例に、各レイヤーがどう連携するか見ていきましょう。

### シナリオ

ユーザーが「日報テンプレート」を使って、新しいノートを作成します。

```
POST /api/notes
{
  "title": "2025年1月23日の日報",
  "template_id": "abc-123",
  "owner_id": "user-456",
  "sections": [
    {"field_id": "field-1", "content": "今日は設計書を書きました"},
    {"field_id": "field-2", "content": "明日はコードレビューします"}
  ]
}
```

### 1. Controller - リクエスト受付

[note_controller.go:73-101](internal/adapter/http/controller/note_controller.go#L73-L101)

```go
func (c *NoteController) Create(ctx echo.Context) error {
    // OpenAPI型でリクエストを受け取る
    var body openapi.ModelsCreateNoteRequest
    if err := ctx.Bind(&body); err != nil {
        return ctx.JSON(http.StatusBadRequest, ...)
    }

    // ドメインDTOに変換
    sections := []port.SectionInput{}
    if body.Sections != nil {
        for _, s := range *body.Sections {
            sections = append(sections, port.SectionInput{
                FieldID: s.FieldId,
                Content: s.Content,
            })
        }
    }

    // UseCaseを呼ぶ
    input, p := c.newIO()
    err := input.Create(ctx.Request().Context(), port.NoteCreateInput{
        Title:      body.Title,
        TemplateID: body.TemplateId.String(),
        OwnerID:    body.OwnerId.String(),
        Sections:   sections,
    })
    if err != nil {
        return handleError(ctx, err)
    }

    // Presenterから結果を取得
    return ctx.JSON(http.StatusOK, p.Note())
}
```

### 2. UseCase - ビジネスロジックの実行

[note_interactor.go:54-90](internal/usecase/note_interactor.go#L54-L90)

```go
func (u *NoteInteractor) Create(ctx context.Context, input port.NoteCreateInput) error {
    // 1. テンプレート取得
    tpl, err := u.templates.Get(ctx, input.TemplateID)
    if err != nil {
        return err  // テンプレートが存在しない
    }

    // 2. ドメインルールで検証
    sections := buildSections("", tpl.Template.Fields, input.Sections)
    if err := note.ValidateNoteForCreate(input.Title, tpl.Template, sections); err != nil {
        return err  // タイトルが空、必須フィールド未入力など
    }

    // 3. トランザクション開始
    var noteID string
    err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        // 4. ドメインサービスでノート構築
        built, err := service.BuildNote(input.Title, tpl.Template, input.OwnerID, sections)
        if err != nil {
            return err
        }

        // 5. ノート保存
        nn, err := u.notes.Create(txCtx, built)
        if err != nil {
            return err
        }
        noteID = nn.ID

        // 6. セクション保存
        sectionsWithID := buildSections(noteID, tpl.Template.Fields, input.Sections)
        if err := note.ValidateSections(tpl.Template.Fields, sectionsWithID); err != nil {
            return err
        }
        return u.notes.ReplaceSections(txCtx, noteID, sectionsWithID)
    })

    if err != nil {
        return err  // トランザクションロールバック
    }

    // 7. 作成されたノートを取得してPresenterに渡す
    n, err := u.notes.Get(ctx, noteID)
    if err != nil {
        return err
    }
    return u.output.PresentNote(ctx, n)
}
```

### 3. Domain - ビジネスルール

[note/logic.go:62-80](internal/domain/note/logic.go#L62-L80)

```go
// ノート作成時の検証
func ValidateNoteForCreate(title string, tpl template.Template, sections []Section) error {
    if strings.TrimSpace(title) == "" {
        return domainerr.ErrTitleRequired  // タイトル必須
    }
    if tpl.ID == "" {
        return domainerr.ErrTemplateOwnerRequired
    }
    if tpl.OwnerID == "" {
        return domainerr.ErrOwnerRequired
    }
    if err := template.ValidateTemplate(tpl); err != nil {
        return err
    }
    if err := ValidateSections(tpl.Fields, sections); err != nil {
        return err  // 必須フィールドチェック
    }
    return nil
}
```

### 4. Domain Service - ノート構築

[service/note_lifecycle.go:9-24](internal/domain/service/note_lifecycle.go#L9-L24)

```go
// ノートを構築
func BuildNote(title string, tpl template.Template, ownerID string, sections []Section) (note.Note, error) {
    if err := note.ValidateNoteForCreate(title, tpl, sections); err != nil {
        return note.Note{}, err
    }
    if ownerID == "" {
        return note.Note{}, domainerr.ErrOwnerRequired
    }
    return note.Note{
        Title:      title,
        TemplateID: tpl.ID,
        OwnerID:    ownerID,
        Status:     note.StatusDraft,  // 必ず下書きで作成
        Sections:   sections,
    }, nil
}
```

### 5. Repository - DB保存

[gateway/db/note_repository.go:130-158](internal/adapter/gateway/db/note_repository.go#L130-L158)

```go
func (r *NoteRepository) Create(ctx context.Context, n note.Note) (*note.Note, error) {
    // 1. string → pgtype.UUIDに変換
    templateID, err := toUUID(n.TemplateID)
    if err != nil {
        return nil, err
    }
    ownerID, err := toUUID(n.OwnerID)
    if err != nil {
        return nil, err
    }

    // 2. sqlcでINSERT
    row, err := queriesForContext(ctx, r.queries).CreateNote(ctx, &sqldb.CreateNoteParams{
        Title:      n.Title,
        TemplateID: templateID,
        OwnerID:    ownerID,
        Status:     string(n.Status),
    })
    if err != nil {
        return nil, err
    }

    // 3. DB行 → ドメインモデルに変換
    return &note.Note{
        ID:         uuidToString(row.ID),
        Title:      row.Title,
        TemplateID: uuidToString(row.TemplateID),
        OwnerID:    uuidToString(row.OwnerID),
        Status:     note.NoteStatus(row.Status),
        CreatedAt:  timestamptzToTime(row.CreatedAt),
        UpdatedAt:  timestamptzToTime(row.UpdatedAt),
    }, nil
}
```

### 6. Presenter - レスポンス整形

[presenter/note_presenter.go:36-41](internal/adapter/http/presenter/note_presenter.go#L36-L41)

```go
func (p *NotePresenter) PresentNote(_ context.Context, n *note.WithMeta) error {
    resp := toNoteResponse(*n)
    p.note = &resp  // 結果を保存
    return nil
}

func toNoteResponse(n note.WithMeta) openapi.ModelsNoteResponse {
    sections := make([]openapi.ModelsSection, 0, len(n.Sections))
    for _, s := range n.Sections {
        sections = append(sections, openapi.ModelsSection{
            Id:         s.Section.ID,
            FieldId:    s.Section.FieldID,
            FieldLabel: s.FieldLabel,
            Content:    s.Section.Content,
            IsRequired: s.IsRequired,
        })
    }

    return openapi.ModelsNoteResponse{
        Id:           n.Note.ID,
        Title:        n.Note.Title,
        TemplateId:   n.Note.TemplateID,
        TemplateName: n.TemplateName,
        OwnerId:      n.Note.OwnerID,
        Owner: openapi.ModelsAccountSummary{
            Id:        n.Note.OwnerID,
            FirstName: n.OwnerFirstName,
            LastName:  n.OwnerLastName,
            Thumbnail: n.OwnerThumbnail,
        },
        Status:    openapi.ModelsNoteStatus(n.Note.Status),
        Sections:  sections,
        CreatedAt: n.Note.CreatedAt,
        UpdatedAt: n.Note.UpdatedAt,
    }
}
```

### まとめ: データの変換フロー

```
HTTP Request (JSON)
    ↓
┌─────────────────────────────────────┐
│ Controller                          │
│ openapi.ModelsCreateNoteRequest     │
│         ↓ 変換                      │
│ port.NoteCreateInput                │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ UseCase                             │
│ ドメインサービス呼び出し              │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ Domain Service                      │
│ note.Note (ドメインモデル)           │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ Repository                          │
│ note.Note                           │
│         ↓ 変換                      │
│ sqldb.CreateNoteParams (sqlc)       │
│         ↓ DB保存                    │
│ sqldb.Note (DB行)                   │
│         ↓ 変換                      │
│ note.Note                           │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ UseCase                             │
│ note.WithMeta 取得                  │
│         ↓                           │
│ Presenterに渡す                     │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│ Presenter                           │
│ note.WithMeta                       │
│         ↓ 変換                      │
│ openapi.ModelsNoteResponse          │
└─────────────────────────────────────┘
    ↓
HTTP Response (JSON)
```

---

## よくある質問と答え

### Q1: なんでこんなに複雑にするの？シンプルじゃダメなの？

**A:** 小さいアプリなら確かにオーバーエンジニアリングです。でも：

```
小規模 (1人、1ヶ月):
┌──────────────┐
│   main.go    │
│  全部ここ！  │ ← これでOK
└──────────────┘

中〜大規模 (チーム、長期保守):
┌──────────────────────────────────┐
│  クリーンアーキテクチャ            │
│  ・テストしやすい                  │
│  ・変更に強い                      │
│  ・チーム開発しやすい               │
│  ・フレームワーク変更できる         │
└──────────────────────────────────┘
```

**このプロジェクトは教材なので、わざと丁寧に分けています。**

---

### Q2: どのレイヤーにテストを書けばいい？

**A:** 全部！でも優先順位は：

```
最優先: ⭐⭐⭐
├─ Domain (ビジネスルール)
│  → ここが壊れたら致命的！
│
├─ UseCase (手順)
│  → ビジネスロジックのテスト
│
次点: ⭐⭐
├─ Presenter (変換)
│  → レスポンス形式が正しいか
│
├─ Gateway (DB変換)
│  → ドメインモデル変換が正しいか
│
最後: ⭐
└─ Controller (HTTPハンドラ)
   → 結合テストやE2Eで確認
```

---

### Q3: Repositoryのインターフェースはどこに置く？

**A:** このプロジェクトでは **Port（`internal/port`）** に置いています。

```
Option A: Portに置く (このプロジェクト)
┌────────────────┐
│   port/        │
│  - NoteRepository interface
│  - TemplateRepository interface
└────────────────┘
         ↑
         │ 実装
┌────────────────┐
│ gateway/db/    │
│  - NoteRepository struct
└────────────────┘

Option B: Domainに置く
┌────────────────┐
│  domain/note/  │
│  - Note entity
│  - NoteRepository interface  ← ここ
└────────────────┘

どっちでもOK！プロジェクトで決めよう。
```

---

### Q4: OpenAPI生成物はどこに置く？

**A:** `internal/adapter/http/generated/openapi/`

**理由:**
- Controllerが使う入力型だから
- Adapterは「外部との変換」を担当する層
- Driverは「配線」だけで、型定義は持たない

---

### Q5: DriverとGatewayの違いは？

```
Driver (internal/driver/)
├─ 役割: 配線、設定、初期化
├─ 例: DB接続、Factory、Initializer
└─ 「準備する人」

Gateway (internal/adapter/gateway/)
├─ 役割: 実際のデータ取得・変換
├─ 例: PostgreSQL Repository、外部API Client
└─ 「実際に動く人」
```

---

### Q6: Factoryパターンを使う理由は？ - 超詳しく解説

**A:** Factoryパターンは一見「面倒」に見えますが、実は**テストと柔軟性**のための超重要パターンです。

#### 🤔 問題: Factory使わないとどうなる？

```go
// ❌ 悪い例: 直接インスタンス化
type NoteController struct {
    usecase *usecase.NoteInteractor  // 具体的な型に依存！
    repo    *db.NoteRepository       // 具体的な型に依存！
}

func NewNoteController(pool *pgxpool.Pool) *NoteController {
    repo := &db.NoteRepository{pool: pool}

    // UseCaseも直接作る
    usecase := &usecase.NoteInteractor{
        notes: repo,
        // ...
    }

    return &NoteController{
        usecase: usecase,
        repo:    repo,
    }
}
```

**何が問題？**

1. **テストが書けない！**
```go
// テストしたいけど...
func TestNoteController_Create(t *testing.T) {
    // ❌ 本物のDBが必要！
    pool, _ := pgxpool.Connect(...)  // DB起動してないと動かない
    ctrl := NewNoteController(pool)

    // ❌ UseCaseをMockに差し替えられない
    // ❌ Repositoryも本物のDBにアクセスする
    // → テスト遅い、環境依存、壊れやすい
}
```

2. **変更に弱い！**
```go
// 将来、MySQLに変えたくなったら...
// → NoteController を全部書き直し！
```

3. **依存が固定される**
```go
// usecase.NoteInteractor の実装を変えたら
// → NoteController も変更が必要
```

---

#### ✅ 解決策: Factory使うとどうなる？

```go
// ✅ 良い例: Factoryパターン
type NoteController struct {
    // 関数（Factory）を保存
    inputFactory    func(...) port.NoteInputPort
    outputFactory   func() *presenter.NotePresenter
    noteRepoFactory func() port.NoteRepository
    txFactory       func() port.TxManager
}

func NewNoteController(
    inputFactory func(...) port.NoteInputPort,
    outputFactory func() *presenter.NotePresenter,
    noteRepoFactory func() port.NoteRepository,
    txFactory func() port.TxManager,
) *NoteController {
    return &NoteController{
        inputFactory:    inputFactory,
        outputFactory:   outputFactory,
        noteRepoFactory: noteRepoFactory,
        txFactory:       txFactory,
    }
}
```

**何が良い？**

1. **テストが簡単！**
```go
// ✅ Mockを簡単に注入できる
func TestNoteController_Create(t *testing.T) {
    // Mock Factory
    mockInputFactory := func(...) port.NoteInputPort {
        return &mock.NoteInteractor{} // テスト用の偽物
    }

    mockOutputFactory := func() *presenter.NotePresenter {
        return presenter.NewNotePresenter()
    }

    mockRepoFactory := func() port.NoteRepository {
        return &mock.NoteRepository{} // テスト用の偽物
    }

    mockTxFactory := func() port.TxManager {
        return &mock.TxManager{} // テスト用の偽物
    }

    // Factoryを差し替えるだけ！
    ctrl := NewNoteController(
        mockInputFactory,
        mockOutputFactory,
        mockRepoFactory,
        mockTxFactory,
    )

    // ✅ DBなしでテストできる
    // ✅ 高速
    // ✅ 安定
}
```

2. **柔軟性が高い！**
```go
// 本番環境
ctrl := NewNoteController(
    factory.NewNoteInputFactory(),      // PostgreSQL版
    factory.NewNotePresenterFactory(),
    factory.NewNoteRepositoryFactory(pool),
    factory.NewTxFactory(pool),
)

// テスト環境
ctrl := NewNoteController(
    testFactory.NewMockNoteInputFactory(),  // Mock版
    testFactory.NewMockPresenterFactory(),
    testFactory.NewMockRepositoryFactory(),
    testFactory.NewMockTxFactory(),
)

// 将来MySQL版が必要なら
ctrl := NewNoteController(
    mysqlFactory.NewNoteInputFactory(),  // MySQL版
    // ...
)
```

---

#### 🎬 実際の動き: リクエストごとに新しいインスタンス

Factoryは「設計図」を保存しておいて、**リクエストごとに新しいインスタンスを作る**ために使います。

```go
func (c *NoteController) Create(ctx echo.Context) error {
    // 1. リクエストを受け取る
    var body openapi.ModelsCreateNoteRequest
    ctx.Bind(&body)

    // 2. Factoryを呼んで、このリクエスト専用のインスタンスを作る
    input, presenter := c.newIO()

    // 3. UseCaseを実行
    err := input.Create(ctx.Request().Context(), ...)

    // 4. Presenterから結果を取得
    return ctx.JSON(200, presenter.Note())
}

func (c *NoteController) newIO() (port.NoteInputPort, *presenter.NotePresenter) {
    // このリクエスト専用のPresenterを作る
    output := c.outputFactory()

    // このリクエスト専用のUseCaseを作る（Presenterを渡す）
    input := c.inputFactory(
        c.noteRepoFactory(),
        c.tplRepoFactory(),
        c.txFactory(),
        output,  // ← Presenterを注入
    )

    return input, output
}
```

**なぜリクエストごとに新しいインスタンス？**

→ **Presenter は結果を保存するから**

```go
type NotePresenter struct {
    note *openapi.ModelsNoteResponse  // ← ここに結果を保存
}

// UseCase がこれを呼ぶ
func (p *NotePresenter) PresentNote(ctx context.Context, n *note.WithMeta) error {
    resp := toNoteResponse(*n)
    p.note = &resp  // ← 結果を保存
    return nil
}

// Controller がこれを呼ぶ
func (p *NotePresenter) Note() *openapi.ModelsNoteResponse {
    return p.note  // ← 保存した結果を返す
}
```

もし**1つのPresenterを複数リクエストで共有**したら...

```
リクエスト1: ノートAを取得 → Presenter.note = ノートA
リクエスト2: ノートBを取得 → Presenter.note = ノートB
リクエスト1: 結果を返す → ❌ ノートBが返ってくる！バグ！
```

**だから、リクエストごとに新しいPresenterを作る必要がある！**

---

#### 🔍 まとめ: Factoryパターンの3つのメリット

```
1️⃣ テストしやすい
   → Mockをサクッと注入

2️⃣ 柔軟性が高い
   → 本番/テスト/MySQL/PostgreSQL...簡単に切り替え

3️⃣ リクエストごとに新しいインスタンス
   → Presenterの結果が混ざらない
```

**結論:**

```
Factory = 「設計図」を保存しておいて、必要なときに新しいインスタンスを作る仕組み

❌ 直接インスタンス化 → テストできない、変更に弱い
✅ Factory → テスト簡単、柔軟、安全
```

---

### Q7: トランザクションはどこで始める？

**A:** **UseCase で開始**します。

```go
func (u *NoteInteractor) Create(ctx context.Context, input port.NoteCreateInput) error {
    // ✅ UseCaseでトランザクション開始
    err := u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        // この中で複数のRepositoryを呼ぶ
        note, err := u.notes.Create(txCtx, ...)
        if err != nil {
            return err  // 自動ロールバック
        }

        err = u.notes.ReplaceSections(txCtx, ...)
        if err != nil {
            return err  // 自動ロールバック
        }

        return nil  // 正常終了 → コミット
    })
    return err
}
```

**なぜ？**
- **ビジネスロジックの境界 = トランザクションの境界**
- Repositoryは「データを取る・保存する」だけに集中

---

### Q8: エラーハンドリングはどうする？

```
ドメインエラー (internal/domain/errors/)
    ↓
ErrNotFound, ErrUnauthorized, ErrTitleRequired...
    ↓
Controller (handleError関数)
    ↓
HTTPステータスコードに変換
    ↓
404 Not Found, 401 Unauthorized, 400 Bad Request
```

**重要:** Gatewayで**DBエラー → ドメインエラー**に変換する！

```go
// ❌ NG: pgx.ErrNoRowsをそのまま返す
if err != nil {
    return nil, err
}

// ✅ OK: ドメインエラーに変換
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domainerr.ErrNotFound
}
```

---

## チェックリスト: コードを書く前に

新しい機能を追加する前に、このチェックリストを確認しよう！

### ✅ Domain層

- [ ] エンティティや値オブジェクトに**フレームワークの型**を使っていないか？
- [ ] ビジネスルール（検証、状態遷移）を**ドメイン層**に書いているか？
- [ ] `import "github.com/..."` などの外部依存がないか？

### ✅ UseCase層

- [ ] Port（インターフェース）経由で依存しているか？
- [ ] 具体的な実装（`&db.Repository{}`）に依存していないか？
- [ ] トランザクション境界を適切に決めているか？
- [ ] 1メソッド = 1仕事になっているか？

### ✅ Controller層

- [ ] HTTPリクエスト → ドメインDTOに変換しているか？
- [ ] UseCaseを呼んで、Presenterから結果を取得しているか？
- [ ] ビジネスロジックを書いていないか？（UseCaseに書く）

### ✅ Presenter層

- [ ] ドメインモデル → OpenAPI型に変換しているか？
- [ ] 結果を内部で保存して、Controllerが取り出せるようにしているか？

### ✅ Gateway層

- [ ] DB行 → ドメインモデルに変換しているか？
- [ ] DBエラー → ドメインエラーに変換しているか？
- [ ] `pgtype.UUID` → `string` などの型変換をしているか？

### ✅ Port層

- [ ] InputPort, OutputPort, Repositoryのインターフェースを定義しているか？
- [ ] 入力DTOを定義しているか？

---

## まとめ: クリーンアーキテクチャの本質

```
┌──────────────────────────────────────────────┐
│  クリーンアーキテクチャの本質                  │
├──────────────────────────────────────────────┤
│                                              │
│  1. ビジネスルールを中心に置く                │
│     → Domainが一番大事                       │
│                                              │
│  2. 外部（DB、HTTP）は交換可能にする          │
│     → Port経由で依存                         │
│                                              │
│  3. 依存の方向を内向きに統一                  │
│     → 外側は内側を知る、逆はダメ              │
│                                              │
│  4. データ構造は各レイヤーで変換              │
│     → OpenAPI型、Domain型、DB型は別物         │
│                                              │
└──────────────────────────────────────────────┘
```

### 最後に

このドキュメントを読んだあなたは、もうクリーンアーキテクチャ初心者ではありません！

**次のステップ:**

1. **実際にコードを読む** - `internal/domain/note/` から始めよう
2. **テストを書いてみる** - Mockを使ってUseCaseをテスト
3. **小さな機能を追加** - 「ノートにタグをつける」とか
4. **リファクタリング** - 既存のコードを改善してみる

**迷ったら:**
- このドキュメントに戻ってくる
- 既存のコードを真似る（`note`を参考にする）
- 3つのルールを思い出す

Happy Coding! 🚀

---

## 付録: ディレクトリ構成の全体像

```
backend-clean/
├── cmd/
│   └── api/
│       └── main.go                      # エントリーポイント
│
├── internal/
│   ├── domain/                          # ❤️ ビジネスルール
│   │   ├── note/
│   │   │   ├── entity.go                # Note, Section
│   │   │   ├── types.go                 # NoteStatus
│   │   │   ├── logic.go                 # 検証ロジック
│   │   │   ├── aggregate.go             # WithMeta
│   │   │   └── *_test.go
│   │   ├── template/
│   │   ├── account/
│   │   ├── service/                     # ドメインサービス
│   │   │   ├── note_lifecycle.go        # BuildNote
│   │   │   ├── status_transition.go     # CanPublish
│   │   │   └── build_sections_from_template.go
│   │   └── errors/
│   │       └── errors.go                # ドメインエラー定義
│   │
│   ├── usecase/                         # 🎯 アプリケーションロジック
│   │   ├── note_interactor.go
│   │   ├── template_interactor.go
│   │   ├── account_interactor.go
│   │   └── mock/
│   │
│   ├── port/                            # 📝 インターフェース
│   │   ├── note_port.go
│   │   ├── template_port.go
│   │   ├── account_port.go
│   │   └── tx.go
│   │
│   ├── adapter/                         # 🔌 外部との接続
│   │   ├── http/
│   │   │   ├── controller/              # HTTPハンドラ
│   │   │   │   ├── note_controller.go
│   │   │   │   ├── template_controller.go
│   │   │   │   ├── account_controller.go
│   │   │   │   ├── server.go            # ルーティング
│   │   │   │   └── mock/
│   │   │   ├── presenter/               # レスポンス変換
│   │   │   │   ├── note_presenter.go
│   │   │   │   ├── template_presenter.go
│   │   │   │   └── account_presenter.go
│   │   │   └── generated/
│   │   │       └── openapi/             # OpenAPI生成物
│   │   │           └── server.gen.go
│   │   └── gateway/
│   │       ├── db/                      # DB Repository
│   │       │   ├── note_repository.go
│   │       │   ├── template_repository.go
│   │       │   ├── account_repository.go
│   │       │   ├── sqlc/                # sqlc生成物
│   │       │   ├── queries/             # SQLクエリ
│   │       │   └── mock/
│   │       └── externalapi/             # 外部API (将来用)
│   │
│   └── driver/                          # 🔧 配線・初期化
│       ├── config/                      # 設定
│       ├── db/                          # DB接続
│       │   ├── pool.go
│       │   └── tx.go
│       ├── factory/                     # Factory関数
│       │   ├── usecase_factory.go
│       │   ├── presenter_factory.go
│       │   ├── repository_factory.go
│       │   └── tx_factory.go
│       └── initializer/
│           └── api/
│               └── initializer.go       # 全体の組み立て
│
├── migrations/                          # DBマイグレーション
├── docs/                                # ドキュメント
└── tests/                               # E2Eテスト (将来用)
```

---

## 付録: 依存関係の図

```
依存の方向: → (知ってOK)

              ┌─────────────┐
              │    Cmd      │ エントリーポイント
              │  (main.go)  │
              └──────┬──────┘
                     │
                     ↓
              ┌─────────────┐
              │   Driver    │ 配線・初期化
              │ (Factory,   │
              │  Config)    │
              └──────┬──────┘
                     │
         ┌───────────┼───────────┐
         ↓           ↓           ↓
    ┌─────────┐ ┌─────────┐ ┌─────────┐
    │Controller│ │ Gateway │ │Presenter│ Adapter層
    └────┬────┘ └────┬────┘ └────┬────┘
         │           │           │
         └───────────┼───────────┘
                     ↓
              ┌─────────────┐
              │   UseCase   │ アプリケーション層
              └──────┬──────┘
                     │
         ┌───────────┼───────────┐
         ↓           ↓           ↓
    ┌─────────┐ ┌─────────┐ ┌─────────┐
    │  Port   │ │ Domain  │ │ Service │ ドメイン層
    │Interface│ │ (Entity)│ │(Domain) │
    └─────────┘ └─────────┘ └─────────┘
         ↑
         │ (実装)
    ┌─────────────────┐
    │ Gateway/Presenter│
    │  (Adapter層)     │
    └─────────────────┘
```

**重要:**
- 矢印は「依存してOK」の方向
- **内側（Domain）は外側を知らない**
- **Port（Interface）で依存を逆転**

---

**このドキュメントがあなたのクリーンアーキテクチャ学習の助けになりますように！** 🎓
