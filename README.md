# Bot
定期的にAPIリクエストを行い、ボトルステータスが`false`であることを検知した場合にDiscordの指定されたチャンネルにボトル状況を投稿するDiscord Botです。`42_done`のリアクションを検知するとボトル交換者と交換までの時間を投稿し、対応したユーザー情報をデータベースへ格納します。全てが終了後、再び次のボトルステータスを問い合わせます。

## シンプルな実行
このプロジェクトを使用するには、Node.jsがシステムにインストールされている必要があります。また、`pnpm`または`npm`を使用して依存関係を管理します。

### インストール手順
1. このリポジトリをクローンします。
    ```bash
    git clone https://github.com/tmyuu/clara-bot.git
    cd clara-bot
    ```

2. `src`ディレクトリ内で依存関係をインストールします。
    ```bash
    cd src
    pnpm install # または npm install
    ```

3. `.env.sample`ファイルをコピーして`.env`ファイルを作成し、適切なトークンを記入します。
    ```bash
    cp .env.sample .env
    ```

4. ボットを実行します。
    ```bash
    node bot.js
    ```

# API
テストと仕様の確認
## 1.データベース
```
$ pwd
   <local_repository_root>/api
$ touch ../db/clara.sqlite3
$ sqlite3 ../db/clara.sqlite3 < ../db/clara.sql
$ sqlite3 ../db/clara.sqlite3
SQLite version 3.45.1
Enter ".help" for usage hints.
sqlite>  INSERT INTO clara(bottle_status, user_name, exchange_status) values('false', 'matsui', 'true');
sqlite> SELECT * FROM clara;
1|false|matsui|true
sqlite>  INSERT INTO clara(bottle_status, user_name, exchange_status) values('true', NULL, 'false');
sqlite> SELECT * FROM clara;
1|false|matsui|true
2|true||false
sqlite> .quit
```

## 2. アプリケーションを起動。
```
$ go run main.go
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.7.2
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:9000
```

## 3.マイコンからAPIへのPOSTリクエスト (ボトル"空"を伝える)
```
$ curl -X POST \
   --url "http://127.0.0.1:9000/empty"
{"message":"Bottle status updated to empty successfully"}
$ sqlite3 ../db/clara.sqlite3
SQLite version 3.45.1
Enter ".help" for usage hints.
sqlite> SELECT * FROM clara;
1|false|matsui|true
2|false||false
sqlite> .exit  
```

## 4. BotからAPIへのGETリクエスト(ボトル"空"の検知)
```
$ curl -X GET 'http://localhost:9000/bottle/2'
{"BottleID":2,"BottleStatus":"false"}
```

## 5. botからAPIへのPOSTリクエスト(user_nameの追加, bottle_idのインクリメントと初期値の代入)
```
$ curl -X POST \
   --url "http://127.0.0.1:9000/adduser" \
   -F "user_name=kamitsui"
{"message":"User name updated successfully"}

$ sqlite3 ../db/clara.sqlite3
SQLite version 3.45.1 2024-01-30 16:01:20
Enter ".help" for usage hints.
sqlite> SELECT * FROM clara;
1|false|matsui|true
2|false|kamitsui|true
3|true||false
sqlite> .exit
```

# 機能

- 定期的にGETリクエストを行い、ボトルステータスの変更を検知します。
- ボトルステータスが`false`になった場合、Discordの任意のチャンネルに投稿します。
- `42_done`のリアクションを検知するとボトル交換者と交換までの時間を投稿し、対応したユーザー情報をPOSTリクエストします。

## 要件
* 水タンクが**"空"**の状態を検知する。
* Botがイベント発生時（水タンクの状態が**"飲料可"**から**"空"**に変わったとき）に特定のDiscordサーバーのチャンネルでメッセージ送信することと、そのチャンネル内でメッセージ検知と書く情報の更新すること。
* Webブラウザからタンクの状態（**"飲料可"** or **"空"**）や現在のタンク在庫数（**最大50本**）やタンク交換したユーザーの情報を確認できる。

## Botが検知するDiscordメッセージとその処理内容
* **"/charged user_name"**のメッセージを検知して、タンクの状態を**"飲料可"**にすること。
* **"/restock n"**のメッセージを検知して、**"校舎在庫"**の値をその値に更新すること。
* **"/reset"**のメッセージを検知して、**"制限数/月"**を**50**に更新すること。


## システムの分担

| No. | **内容**                     | 機器         | 担当                        |
| :-- | :------------------------- | :--------- | :------------------------ |
| 1.  | **水槽水位検知センサー**             | 電子回路       | kamitsui                  |
| 2.  | **WiFi 機能を備えたマイクロコントローラー** | マイコン       | kamitsui                  |
| 3.  | **サーバー バックエンド**            | Jetson \*1 | ymatui, kamitsui, ryozuno |
| 4.  | **データベース**                 | Jetson \*1 | ?                         |
| 5.  | **Discord ボット**            | Jetson \*1 | ymatsui                   |
| 6.  | **Web インターフェイス**           | Jetson \*1 | ryozno                    |

 \*1　実験時はJetsonで、運用時は42Tokyoサーバーを使えたら嬉しい。

## Mermaid
```mermaid
---
title: Develop system diagram (A)
---

%% graph %%
graph LR

%% subgraph %%
subgraph "Water Server"
    sensor(Sensor)
    microcontroller[MCU]
end

subgraph "Jetson"
  subgraph "Main"
    subgraph "Container_01"
      serverBackend[Backend]
    end
    subgraph "Container_02"
      discordBot[Bot]
    end
    subgraph "Container_03"
      database[DB]
    end
  end
  subgraph "Container_04"
    webInterface[Web_Server]
  end
end

subgraph "Discord_Server"
    discordInterface[Aqua_Clara_CH]
end

%% %% Node to Node
%% sensor -> micon
sensor -->
  |ON/OFF
  *1| microcontroller

%% mcu -> Backend
microcontroller -->
  |HTTP *2| serverBackend

%% Back -> DiscordBot
serverBackend -->
  |*3| discordBot
discordBot -->
  |*4| serverBackend

%% DiscordBot <-> Discord Server
discordBot -->
  |*5| discordInterface
discordInterface -->
  |*6| discordBot

%% Backend -> DataBase
serverBackend -->
  |*7| database
database -->
  |*8| serverBackend

%% Web Server <-> Backend
webInterface -->
  |*9| serverBackend
serverBackend -->
  |*10| webInterface
%% %% -----------
```

## 貢献

プロジェクトへの貢献に興味がある方は、プルリクエストを送るか、[x](https://x.com/livspect)から連絡してください。


