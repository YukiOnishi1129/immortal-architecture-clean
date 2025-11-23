# はじめてのテスト設計ガイド（backend-clean）

このドキュメントは「駆け出しエンジニアでも迷わず読める」ことを目的に、テストの考え方と実際のコード配置をまとめたものです。Clean Architecture と DDD の設計方針に沿いつつ、実践的なテストのやり方を手順で示します。

## 基本方針
- **レイヤーごとに責務を分けてテスト**する。外側のレイヤーは内側を直接触らず、ポート（インターフェース）越しにモックを差し込む。
- **C1（主要な分岐網羅）**を目安に、正常系と代表的な異常系を必ず入れる。
- **外部依存はモック**で置き換える。DB・HTTPクライアント・TxManager などはモック化し、テストDBは使わない（Gatewayもモックユニットで担保）。
- **テーブル駆動テスト**を使い、`[Success] ...` / `[Fail] ...` の名前でケースを明示する。
- **エラーメッセージ**も検証する。HTTPレスポンスはステータスコードだけでなくボディに期待文字列が含まれるかを確認する。

## レイヤー別の考え方
- **Domain**: 純粋関数が多いので直接テスト。値オブジェクトのバリデーション、状態遷移などのルールを正常/異常で確認する。
- **UseCase**: ポートを gomock（もしくは手書きスタブ）でモック。リポジトリ/TxManager/OutputPort の呼び出しが期待どおりかを確認する。
- **Controller/Presenter**: Echo のコンテキストをフェイクし、Bind失敗→400、owner不足→403、NotFound→404などのマッピングをテーブル駆動で確認。Presenterはドメイン→OpenAPI DTOへの詰め替えをチェック。
- **Gateway**: SQLC/pgx をモックし、UUID変換エラー、ErrNoRows→ErrNotFound へのマッピング、ドメイン型への詰め替えをテーブル駆動で確認。実DBは使わない。
- **Driver/Initializer**: 配線のみ。ユニットではスモークテスト（nil依存でもパニックしないこと）を入れる。実行時はビルド/起動で検知。

## ディレクトリとファイル配置（抜粋）
- テストは同一ディレクトリに置き、`*_test.go` で管理。
- モックはレイヤー配下に置く。
  - `internal/usecase/mock/*` : UseCase用モック（gomock生成物）
  - `internal/adapter/gateway/db/mock/*` : Gateway用DBTXモック（手書き）
  - `internal/adapter/http/controller/mock/*` : Controller用スタブ
- lint は `*_test.go` と mock ディレクトリを除外済み（`.golangci.yml`）。

## 何をテーブルに入れるか（例）
- ドメイン: 正常系＋バリデーションNG、状態遷移NGなど。
- UseCase: 正常系＋リポジトリエラー、Txエラー、Unauthorized/NotFoundなど。
- Controller: Bind失敗、owner不足、UseCaseエラー（Forbidden/NotFound）、正常系。
- Gateway: UUID不正、ErrNoRows→ErrNotFound、Query/Execエラー、置換系（ReplaceSections/ReplaceFields）でのエラー。

## 実行方法
- テスト全体: `GOCACHE=/tmp/gocache go test ./...`
- （lintはネットワーク制限下で失敗する場合あり。mock と *_test.go は除外設定済み）

## ゴールイメージ
- 新しく分岐やドメインルールを追加したら、同じディレクトリの `*_test.go` にテーブル駆動でケースを足す。
- 「正常系＋代表的な異常系」を常にセットで書き、ステータスやメッセージを確認する。
- 外部依存に触る処理はモックで閉じ込め、テストDBなしで高速に回る状態を保つ。
