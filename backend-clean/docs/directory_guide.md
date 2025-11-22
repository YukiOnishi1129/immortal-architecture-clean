# ディレクトリ構成ガイド（ポート集約スタイル）

Clean Architecture をこのプロジェクトでどう落とし込むかを、図解＋箇条書きでまとめます。ポートは `internal/port` に集約し、ドメインはパッケージ分割します。

```
backend-clean/
├─ cmd/api/main.go                 # エントリーポイント（driverを呼び出す）
└─ internal/
   ├─ port/                        # 境界の契約を集約（ドメイン別ファイル）
   │   ├─ account_port.go          # Account Input/Output/Repository/DTO
   │   ├─ note_port.go             # Note Input/Output/Repository/DTO
   │   ├─ template_port.go         # Template Input/Output/Repository/DTO
   │   └─ tx_port.go               # TxManager 契約
   ├─ domain/                      # 純粋なモデルとルール（ドメイン単位のパッケージ）
   │   ├─ account/                 # package account
   │   │   ├─ entity.go            # エンティティ/VO
   │   │   └─ logic.go             # ドメインロジック・検証
   │   ├─ note/                    # package note
   │   │   ├─ entity.go
   │   │   └─ logic.go
   │   └─ template/                # package template
   │       ├─ entity.go
   │       └─ logic.go
   │   └─ service/                 # 複数ドメインをまたぐドメインサービス（必要なら、サービスごとにファイル分割）
   ├─ usecase/                     # Interactor（ポート実装、FW/DBを知らない）
   │   ├─ account_interactor.go
   │   ├─ note_interactor.go
   │   └─ template_interactor.go
   ├─ adapter/
   │   ├─ http/
   │   │   ├─ controller/          # ServerInterface実装（入力アダプタ）
   │   │   │   ├─ account_controller.go
   │   │   │   ├─ note_controller.go
   │   │   │   ├─ template_controller.go
   │   │   │   └─ helper.go
   │   │   ├─ presenter/           # 必要なら出力整形
   │   │   └─ generated/openapi/   # oapi-codegen 生成物
   │   └─ gateway/
   │       ├─ db/                  # sqlc生成物＋DBリポジトリ実装
   │       └─ externalapi/         # 外部API実装（必要に応じて）
   └─ driver/                      # 配線・初期化（ビジネスロジックなし）
       ├─ db/                      # Pool/TX などDB関連初期化
       ├─ config/                  # 環境変数など設定読み込み
       └─ initializer.go           # Echo起動・ハンドラ登録など
```

## 矢印と責務
- ポートは `internal/port` に一本化。UseCaseとAdapterはここに依存する。
- Domain は純粋なモデル/ロジックのみで、ポートやFWを知らない。
- UseCase はポート経由で Domain を動かし、Presenter/Gateway/Tx を呼ぶ。
- Adapter はポートを実装する側（HTTP Controller/Presenter、DB Repository など）。
- Driver は Echo/DB/Config の初期化と配線だけを担当。ビジネスロジックは持たない。

この構成で「どこに何を書くか」を迷わず進められます。***
