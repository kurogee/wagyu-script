v1.2.0-beta (2024-08-06)
---
 - README.mdを更新しました
 - Wagyu-scriptのパッケージの反映忘れが発覚したため、v1.1.0-betaより前のバージョンが非推奨となりました
 - docs/sharp-function.mdを修正しました
 - docs/common-function.mdを修正しました
 - docs/get-started.mdを修正しました
 - README.mdを更新しました
 - dateパッケージにシャープ関数を追加しました
 - dateパッケージに曜日を取得する関数を追加しました
 - 変数を複数宣言する、vars関数を追加しました
 - **docsに英語版を追加しました / Added English version to docs**
 - ↑に伴い、docs/jaとdocs/enを作成しました
 - ↑の更新がありましたが、今後のCHANGELOGでも``docs/_____.mdを修正しました``という表記を使います

v1.1.0-beta (2024-08-05)
---
 - バージョン表記の仕組みを変更しました（beta.1→betaというようにbeta以降の数字を廃止など）
 - 返り値関係のバグを修正しました
 - swap関数を追加しました
 - docs/get-started.mdとdocs/common-function.mdを修正しました
 - サンプルコードを置いておくdocs/sample-code.mdを作成しました
 - README.mdを更新しました
 - ロゴを作成しました
 - arrayパッケージで、半角スペースを不適切に処理していたのを修正しました
 - arrayパッケージにjoin関数を追加しました
 - 変数を宣言する時の「value」や「array」を省略できるようにしました

v1.0.3-beta.1 (2024-08-05)
---
 - シャープ関数を自分で定義できるようになりました
 - シャープ関数にlen関数、repeat関数を追加しました
 - シャープ関数の入れ子が可能になりました
 - ↑に伴い、docs/sharp-function.mdを修正しました

v1.0.2-beta.1 (2024-08-04)
---
 - regexパッケージを追加しました（関数とシャープ関数の両方に）
 - シャープ関数にvar関数を追加しました
 - match文を追加しました
 - バッククオートを用いて文字列を表記できるようになりました
 - each関数に渡す引数を変更しました
 - docs内に新しい機能・関数の説明を追加しました
 - get-started.mdの説明を修正しました

v1.0.1-beta.3 (2024-08-04)
---
 - CHANGELOG.mdの表記を変更しました
 - docs/get-started.mdを修正しました
 - ↑each関数の説明を追加しました
 - each関数で文字列のクオートが取れずに出力されてしまう問題を修正しました

v1.0.1-beta.2 (2024-08-03)
---
 - パッケージが上手く反映されない問題を修正しました
 - each関数を追加しました

v1.0.1-beta.1 (2024-08-03)
---
 - docs/get-started.mdを追加しました
 - docs/common-function.mdを追加しました
 - docs/sharp-function.mdを追加しました
 - README.mdを更新しました
 - CHANGELOG.mdを追加しました
 - CHANGELOG.mdを記録し始めました
 - arrayパッケージの引数の一部が変数に対応していなかったので修正しました
 - fileパッケージで、\nや\tなどのエスケープシーケンスを正しく処理できるようにしました
 - mathパッケージで、数値の切り下げ、切り上げ、四捨五入を行う関数を追加しました