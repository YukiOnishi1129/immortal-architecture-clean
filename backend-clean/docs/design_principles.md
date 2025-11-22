# Backend Clean API 設計ガイド（初心者向け）

アンチパターン版（bad-api）の「つらさ」を解消するため、ここでは **Clean Architecture + DDD** を「現場の手順」としてまとめます。初めて読む人でもイメージしやすいよう、たとえ話を交えて説明します。

---

## まず全体像：建物のイメージ
- **建物の芯 = Domain**  
  変わらない耐震構造。ビジネスの概念（Note, Template など）とルールだけを書く。外の素材（FW/DB 型）は持ち込まない。
- **廊下 = UseCase**  
  人の動線。「誰が何をするか」を一つのメソッドにまとめ、必要なら鍵（トランザクション）を管理する。ここで「扉の形（Port）」を決める。
- **扉 = Adapter**  
  - Controller: HTTP から入ってきたものを廊下用の形に詰め替える。
  - Presenter: 廊下の結果を外に見せる形に詰め替える。
  - Gateway: DB や外部 API に出入りする扉。廊下が決めた形に合わせる。
- **外壁・設備 = Driver**  
  Echo や DB 接続、設定、ログなどの初期化。内側の設計を変えずに貼り替えられる（Infrastructure というディレクトリは作らず driver にまとめる）。
- **矢印は内向きのみ**  
  外壁 → 扉 → 廊下 → 芯 の一方通行。芯（Domain）は扉や外壁の存在を知らない。

---

## 目標（なぜこうするか）
- **変更に強くする**: FW/DB/ORM を差し替えても、ドメインとユースケースのコードをほぼ触らずに済む。
- **テストしやすくする**: 内側はモックで閉じ、外側だけ統合/E2E テストにする。
- **責務を迷わない**: 「どこに何を書くか」のルールを固定して再現性を上げる。

---

## ディレクトリのガイド
- `cmd/api/main.go` … エントリーポイント（配線だけ）
- `internal/domain` … エンティティ、値オブジェクト、ドメインサービス、ドメインエラー
- `internal/usecase` … ユースケース（Input/Output DTO と Port をここに置く）
- `internal/adapter/controller` … HTTP/CLI などの入力アダプタ
- `internal/adapter/presenter` … 出力アダプタ（API レスポンスを組み立てる）
- `internal/adapter/gateway` … DB/外部 API への実装（Repository/Gateway 実体）
- `internal/driver` … FW/DB/設定/ログなどの初期化（Infrastructure 用の別ディレクトリは作らない）

ポイント: **Port は usecase 配下に置く。Adapter 側には置かない。**

---

## ルールセット（初心者がまず覚える7つ）
1) **依存は内向きだけ**  
   どの層も「内側のインターフェース」にだけ依存する。内側は外側を知らない。
2) **境界ははっきり書く**  
   入力ポート/出力ポート、DTO/モデルの変換場所を固定する（Controller/Presenter/Gateway に限定）。
3) **ドメインは純粋に**  
   FW/DB 型を持ち込まない。値オブジェクトで不変条件を守る。
4) **ユースケースは1メソッド**  
   ユースケースごとにメソッドを分け、ここでトランザクション境界を決める。
5) **共通部品も注入**  
   ログ/設定/時刻/ID 生成もインターフェース越しに注入。new で直接作らない。
6) **テストはモックで閉じる**  
   内側はモックでテストできる設計にする。DB や外部 API 依存は統合テストで確認。
7) **データ変換の道を固定**  
   OpenAPI DTO ↔ UseCase DTO ↔ Domain の順で詰め替え、逆も同じ道を戻る。

---

## レイヤ別の役割としてほしいこと
- **Domain**: ビジネスの意味だけを書く。例: Note のステータスは Draft/Publish だけ、Template の Field は順番重複不可、など。
- **UseCase**: ユースケースの手順書。Input DTO を受けてドメインを動かし、Output DTO を返す。Port（Repository/Gateway/Presenter/TxManager）にだけ依存。
- **Controller**: HTTP の形を UseCase DTO に詰め替える。追加バリデーションもここ。
- **Presenter**: UseCase の結果を API レスポンスに詰め替える（HTTP ステータス決定もここ）。
- **Gateway**: DB/外部 API との橋渡し。Domain モデルを入出力に使う。
- **Driver**: Echo や DB の配線。内側のコードを触らずに差し替えられるようにする（Infrastructure 用の別ディレクトリは持たない）。

---

## トランザクションの扱い
- UseCase が `RunInTx(ctx, func(ctx) error)` を呼び、境界を明示する。
- Tx の実装は Gateway/Driver に閉じ込める。UseCase は「実行してほしい関数」を渡すだけ。
- ネストは禁止（必要なら呼び出し階層を設計で整理）。

---

## エラーの扱い
- ドメインエラーを型で表現（例: `TemplateNotFound`, `TemplateInUse`, `InvalidStatusTransition`）。  
  Presenter が HTTP ステータスに写像する。
- 予期せぬエラーはラップしてログし、クライアントには安全なメッセージだけ返す。

---

## テストの考え方
- **ユニットテスト**: UseCase は Repository/TxManager/Gateway/Presenter をモック。Domain は純粋関数として直接テスト。
- **統合テスト**: testcontainers などで実 DB を起動し、Gateway 実装とマイグレーションを検証。
- **E2E**: 必要な主要フローだけを API 経由で確認。

---

## 変換の具体例（Template 作成の流れ）
1. HTTP リクエスト（OpenAPI DTO）  
   → Controller が UseCase Input DTO に詰め替え  
2. UseCase がドメインを操作  
   → Repository Port 経由で保存  
3. UseCase Output DTO  
   → Presenter が API レスポンスに詰め替え  
4. HTTP レスポンスとして返す

---

## まとめ
このガイドのゴールは「どこに何を書くかを迷わないこと」と「変更とテストに強い状態を保つこと」です。建物のたとえを意識しながら、アンチパターン版で感じた痛み（境界なし・手書きトランザクション・型混在・テスト不能）をここで順に潰していきましょう。
