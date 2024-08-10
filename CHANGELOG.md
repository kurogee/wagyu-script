v1.6.1-beta (2024-08-10)
---
 - 辞書型(dictパッケージ)の変数に関する処理の不具合を修正しました
 - docs/sharp-function.mdを修正しました

v1.6.0-beta (2024-08-10)
---
 - Wagyu-scriptシステムのsplit、take_off_quotation、divide_split関数を1つのファイルにまとめて、管理をしやすくしました
 - 辞書型を追加しました
 - dictパッケージを追加しました
 - docs/samples.mdを編集しました
 - docs/get-started.mdを編集しました
 - docs/sharp-function.mdを編集しました
 - docs/common-function.mdを編集しました

v1.5.0-beta (2024-08-09)
---
 - arrayパッケージのseach関数で正しく値が処理されない問題を修正しました
 - arrayパッケージのsplitなどの関数で配列が適切に処理されない問題を修正しました
 - arrayパッケージにシャープ関数split、joinを追加しました
 - stringパッケージにシャープ関数include、substrを追加しました
 - 括弧内のエスケープシーケンスを正しく処理できるようにしました
 - docs/sharp-function.mdを更新しました
 - README.mdを更新しました

v1.4.3-beta (2024-08-08)
---
 - match文で、「_」を値に置き換えるということができていなかったので修正しました

v1.4.2-beta (2024-08-08)
---
 - dateパッケージのadd関数、sub関数で、引数がうまく渡らない問題を修正しました
 - dateパッケージのadd関数、sub関数で、引数が変数に完全に対応しました

v1.4.1-beta (2024-08-07)
---
 - README.mdが間違えて更新されていなかったので修正しました

v1.4.0-beta (2024-08-07)
---
 - シャープ関数のfrom関数の処理速度が向上しました
 - stringパッケージにシャープ関数replaceを追加しました
 - arrayパッケージにシャープ関数searchを追加しました
 - 配列処理を行う際に、文字列のスペースと配列を区切るためのスペースが区別されずに処理されていたのを修正しました
 - シャープ関数arrayAtで、数値の配列と文字列の配列を区別できるようになりました
 - docs/sharp-function.mdを更新しました
 - docs/common-function.mdを更新しました
 - README.mdに「Credits」と「License」項目を追加しました

v1.3.1-beta (2024-08-07)
---
 - シャープ関数の返り値がうまく取得できない問題を修正しました

v1.3.0-beta (2024-08-07)
---
 - each関数に渡す引数が、シャープ関数に対応しました
 - case文につく値が、変数に対応しました
 - varsが、シャープ関数に対応しました
 - delete関数が、複数の変数を削除できるようになりました
 - シャープ関数atを追加しました
 - シャープ関数allを追加しました
 - シャープ関数anyを追加しました
 - regexパッケージにシャープ関数findを追加しました
 - プログラム側が、渡された引数が変数か値かを判別できるようになりました
 - ↑に伴い、**予期せぬバグが起こるかもしれません**
 - docs/common-function.mdを修正しました
 - docs/sharp-function.mdを修正しました
 - docs/get-started.mdを修正しました

v1.2.0-beta (2024-08-06)
---
 - README.mdを更新しました
 - Wagyu-scriptのパッケージの反映忘れが発覚したため、v1.1.0-betaより前のバージョンが**非推奨**となりました
 - docs/sharp-function.mdを修正しました
 - docs/common-function.mdを修正しました
 - docs/get-started.mdを修正しました
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

##### ※これより下のバージョンは非推奨です

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