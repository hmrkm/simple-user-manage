# Simple-user-manage

[![codecov](https://codecov.io/gh/hmrkm/simple-user-manage/branch/main/graph/badge.svg?token=LE4923URW1)](https://codecov.io/gh/hmrkm/simple-user-manage)

シンプルなユーザー管理API

## 必要なもの

- Docker Compose

## インストール

1. `.env.sample`をコピーして`.env`を作成
2. `.env`の内容を修正
3. `auth/.env.sample`をコピーして`auth/.env`を作成
4. `auth/.env`の内容を修正
5. `docker-compose up -d`
6. DBに`app/docs/migration.sql`の内容を反映
7. `app/docs/insert_test_user.sql`でテストデータを追加