# テスト方針（Unit / C1 ホワイトボックス）

- 目的: 各レイヤーの責務を単体で担保する。外部環境（DB・認証等）には依存しない。
- カバレッジ: C1（分岐網羅）を目安に、正常系と主要な異常系を押さえる。
- ツール: `testing` 標準パッケージ + `gomock`。テスト DB は使わない。
- スタイル: テーブル駆動、Given/When/Then コメント、モックの期待値は最小限かつ明示。

## レイヤー別の観点とモック対象
- **Domain（Entity/Logic）**: 外部依存なし。純粋関数を直接テスト。境界値・不変条件を重点。
- **UseCase**: Ports を gomock でモック。リポジトリ/TxManager/OutputPort などの呼び出し順・引数を確認。副作用の有無と戻り値を検証。
- **Adapter (Gateway/Controller/Presenter)**:
  - Gateway: SQLC クエリはモック化し、入力/出力マッピングを確認（ドメイン型への詰め替え含む）。
  - Controller: Echo の Context は最小のフェイク/モックで Bind/JSON 応答を確認（Body 変換・パラメータ検証・ステータスコード）。
  - Presenter: ドメイン→OpenAPI DTO への変換を検証（フィールド欠落・ゼロ値の扱い）。
- **Driver/Initializer**: 依存組み立てのみ。ユニットテスト対象外（統合テストで確認）。

## 典型的なテスト構成
- パス: `internal/<layer>/*_test.go`
- テーブル駆動: `{name, input, mocks, wantErr, want}` を定義
- gomock: `mockgen` で Port/Gateway/TxManager/OutputPort のモックを生成し、`EXPECT()` で呼び出し回数と引数を検証
- エラーパス: 想定される代表的な異常を 1 ケース以上入れる（例: バリデーション NG, リポジトリエラー, トランザクション失敗）

## カバレッジ観点（例）
- **UseCase/TemplateInteractor**: 
  - List: フィルタ無し/所有者指定あり/リポジトリエラー
  - Create: 正常系（Fields あり/なし）、バリデーション NG、Tx 内で ReplaceFields 失敗
  - Update: オーナー不一致、フィールド空、テンプレ使用中エラー、正常更新
  - Delete: オーナー不一致、TemplateInUse、正常削除
- **UseCase/NoteInteractor**:
  - Create: テンプレート存在チェック、セクション自動生成、ValidateSections NG
  - Update: 所有者不一致、セクション置換（既存/追加）、ValidateSections NG
  - ChangeStatus: Draft→Publish 正常、Publish→Draft 正常、無効遷移、オーナー不一致
  - Delete: 所有者不一致、正常削除
- **Gateway (template/note/account)**:
  - DB 行→ドメインへのマッピング（nullable, time 変換）
  - エラー伝搬（pgx.ErrNoRows→domain.ErrNotFound など）
- **Controller (note/template/account)**:
  - Bind 失敗で 400
  - ownerId 不足で 403
  - UseCase エラーのマッピング（NotFound→404, Unauthorized→403 など）
  - 正常時のレスポンスステータス/ボディ

## モック生成のメモ
- Port/Repository/TxManager/OutputPort は `mockgen` で生成し `internal/mock/` 配下に置く想定（例: `mock_port`, `mock_usecase`）。
- 生成コマンド例（参考）:
  - `mockgen -source=internal/port/note_port.go -destination=internal/mock/port/mock_note_port.go -package=portmock`
  - `mockgen -source=internal/port/tx_port.go -destination=internal/mock/port/mock_tx_port.go -package=portmock`
  - `mockgen -source=internal/adapter/gateway/db/sqlc/queries.go -destination=internal/mock/gateway/mock_queries.go -package=gatewaymock`

## 運用
- テスト実行: `GOCACHE=/tmp/gocache go test ./internal/...`
- CI: lint/test のみ（DB を使わないため追加セットアップ不要）

## gomock 配置ルールと残タスク
- モックはレイヤー配下に置く: 例 `internal/usecase/mock/*`（パッケージ名は `mockusecase` など衝突しないものに統一）。
- Controller/Gateway用のモックが必要な場合は、EchoやSQLCクエリに薄いインターフェースを定義してそこに対して `mockgen` を当てる。
- テストケース名は `[Success] ...` / `[Fail] ...` を先頭に付けて正常/異常を明示する。
- 進め方（残りの主なタスク）:
  - Domain: 新規ドメインが増えたら logic/aggregate/service を同様にテーブル駆動で追加。
  - UseCase: NoteInteractor/AccountInteractor に gomock テストを追加（正常・バリデーションNG・リポジトリエラー・Txエラー）。
  - Gateway: SQLCクエリをモックし、DB行→ドメイン変換と ErrNoRows→ErrNotFound の変換を検証。
  - Controller: Bind失敗/ownerId不足/UseCaseエラー→HTTPコードのマッピングを確認するフェイク Echo Context テストを追加。

この方針に沿って各レイヤーのユニットテストを追加していきます。
