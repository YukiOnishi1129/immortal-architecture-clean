# ポートとアダプターの設計メモ

「なぜポートがいるの？」「Controller / Presenter / Gateway / Driver って何？」を初学者向けにまとめました。

## ざっくり全体像
- **Controller（入り口）**: HTTP リクエストを受け取り、UseCase を呼ぶ係。Echo など FW 依存はここに閉じ込める。
- **UseCase（ビジネス手順）**: ドメインを動かし、結果を Presenter に渡す。FW や DB の型を知らない。
- **Presenter（出口の整形係）**: UseCase の結果をレスポンス用に詰め替える（HTTP ステータスや JSON など）。
- **Gateway（外部と話す係）**: DB や外部 API とやりとり。Repository 実装や外部クライアントなど。
- **Driver（配線係）**: Echo 起動、ルーティング登録、DB 接続や TxManager の生成。ビジネスロジックは持たない。

矢印は外→内だけ。内側は外側を知らない。

## ポートを作る理由（3行で）
1. **混ざらない**: FW/ORM の型が内側に入ってこない。
2. **差し替えやすい**: HTTP→CLI など入口変更や、DB 実装入替えがしやすい。
3. **テストしやすい**: モックに差し替えて、重い外部依存なしでユニットテストできる。

## HTTP の流れ（例）
1. Controller: HTTP を受けて、UseCase InputPort を呼ぶ。
2. UseCase: ドメインを動かし、結果を OutputPort(Presenter) に渡す。Tx や Repository もポート経由。
3. Presenter: ドメイン結果をレスポンス用 DTO / HTTP ステータスに変換。
4. Controller: Presenter からもらった形でレスポンスを返す（Presenter が直接書き出す設計でも可）。

## このプロジェクトの置き場所・名前
- **ポート（契約を集約）**: `internal/port/` にドメインごとファイル分割  
  - 例: `account_port.go`, `note_port.go`, `template_port.go`, `tx_port.go`  
  - 中身: InputPort / OutputPort / Repository インターフェース + UseCase DTO
- **UseCase（Interactor）**: `internal/usecase/account_interactor.go` など。ポートを実装するだけで、FW/DBを知らない。
- **Domain（純粋なモデルとルール）**: ドメイン単位のパッケージに分割  
  - 例: `internal/domain/account/{entity.go,logic.go}` / `internal/domain/note/...` / `internal/domain/template/...`  
  - 複数ドメインにまたがるドメインサービスは別ディレクトリで区別する。
- **Adapter**:  
  - HTTP: `internal/adapter/http/controller/*`, `internal/adapter/http/presenter/*`, 生成物は `internal/adapter/http/generated/openapi`  
  - Gateway: `internal/adapter/gateway/db`（sqlc生成物＋リポジトリ実装）、`.../externalapi` など外部リソースごとにサブディレクトリ
- **Driver（配線・初期化）**: `internal/driver/db`（接続/TX）, `internal/driver/config`, `internal/driver/initializer.go`（Echo 起動とハンドラ登録）など

## よくある疑問
- **UseCase を直接呼べばいい？**  
  ポートがあると入口を変えたり、テストでモックしやすい。境界も型でハッキリする。
- **OutputPort は何をする？**  
  UseCase が HTTP を知らないまま結果を渡す窓口。Presenter が FW 依存の整形を担当。
- **TxManager をポートにする意味は？**  
  トランザクション管理の契約を内側が握り、実装を外側（driver）に追い出すことで、テストや差し替えを容易にする。

設計方針（`docs/design_principles.md`）とセットで読んで、Controller/Presenter/Gateway/Driver を適切な場所に置いてください。
