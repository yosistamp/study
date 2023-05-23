# Global Azure

#GloalAzure #jazug

- 日付：2023/5/13
- 参加形態：オンライン（マイクロソフト品川オフィスでもやっている）

## Virtual EC WebSite Implemented by multi lang run on Azure Spring Apps.

### Azure Spring Apps

C# Python Node.js Go などを動かせる

### サンプルアプリ

https://github.com/Azure-Samples/acme-fitness-store

- 認証
  - AzureADと連携
- Monitoring
  - Javaはソースコードがあれば自動でApplicationInsightが使える。ほかの言語は手動設定。（使えはする）
- シークレット
  - AzureKeyVault

### AzureOpenAIでの問い合わせ

ChatGPTに事前にサイトの問い合わせの身になるように設定しておけば、関係のない問い合わせには答えないようにできる。

### AzureSpringApps と Kubernetes

単純な価格はPaaSのほうが高い。ただし人件費を含めて考えたほうがいい。

- 単純なサービス公開はAppService
- サービスが増えたり、サービス間連携が増えたりしたらAzureSpringCloud

AzureContainerApps

Dockerのイメージを作って、若干Kubernetesを意識する

AzureSpringAppsはコンテナ作成も不要

### AzureSpringApps

プランはbasic、standard、Enterprise

java以外はEnterpriseプラン

#### アプリケーション作成

- az spring create
  - クラスタ作成
- az spring app create
  - ServiceやDeploymentの準備（アプリはない）

ロールの割り当てをすればブラウザからコンテナにログインできる

Blue/Greenもできる az spring deployment create (Deploymentを分ける)

### みんなの反応

AzureContainerAppsとかAppServiceとの使い分けが、やっぱりしっくり来てない感じ。
あと、Springに限らないのにSpringAppsは名前が悪いとの意見。（そりゃそうだ）

## Pulumi de Azure IaC (C+Dブース)

プログラミング言語でIaCができる（CDKみたいなもの）

AWS/Azure/GCP/KubernetesのAPIサポート率100%

## Azure Policyとガバナンスのおはなし (C+Dブース)

財務/ビジネス→ロール→リソースグループにタグ付け

### AzurePolicy

リソースがビジネスルールに準拠しているか定義
親スコープは子スコープを継承

VMのサイズ限定
タグの必須化
カスタムポリシー

### RBACとの違い

RBACは権限定義、やっていいことの権限付与

### CAF（CloudAdpointFramework）

ベストプラクティス、ドキュメント、ツール

内容はaz2.jpg

#### AzureBluePrint

- 決められたリソースや環境を一括でデプロイ
- ガバナンス系のテンプレートもいくつか用意されている
- ずっとPreview

組織が大きくなるとどこで設定されたポリシーかわかりづらくなる

#### AzureCovernanceVisualizer

GithubActionsで自動化できる


## Azure OpenAI Service + Semantic Kernel (C#) (C+Dブース)

https://youtu.be/tFgqdHKsOME

### ChatGPT

アシスタントのセットアップでシステムメッセージを入れておくと前提条件になる。

- Azureと本家で価格は同じ
- SLAはAzureのみある
- セキュリティもAzureはVNetのみからのアクセスやマネージドIDで管理できる

### SemanticKarnel

- スキルの実行順を管理
- スキルのプロンプトを管理
- メモリ管理
- 最大文字を超えるトークン管理
- 最新情報を外部から取得
- モジュール化

## 突如登場した Azure Developer CLI でなにができるのか？検証してみる (Bブース)