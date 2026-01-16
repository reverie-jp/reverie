### Getting Started

> [!TIP]
> 開発には [Visual Studio Code](https://code.visualstudio.com/) と [Docker](https://www.docker.com/) を使用します。あらかじめインストールをお願いします！

1. リポジトリをクローン

```shell
git clone https://github.com/reverie-jp/reverie.git
```

2. 環境変数をコピー

```shell
cp .env.development .env.development.local
```

3. Visual Studio Code に [Dev Containers 拡張機能](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) をインストールし、「Reopen in Container」を選択

4. API サーバーを起動

```shell
make dev-up
```
