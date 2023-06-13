const path = require("path");
module.exports = {
	// モジュールバンドルを行う起点となるファイルの指定
	// 指定できる値としては、ファイル名の文字列や、それを並べた配列やオブジェクト
	// 下記はオブジェクトとして指定した例
	mode: "development",
	devtool: "inline-source-map", // inline- にすると、bundle にソースマップが含まれるため、ファイルサイズが大きくなる
	entry: {
		main: [
			"./typescript/index.tsx"
		]
	},
	output: {
		// モジュールバンドルを行った結果を出力する場所やファイル名の指定
		// "__dirname"はこのファイルが存在するディレクトリを表すnode.jsで定義済みの定数
		path: __dirname,
		filename: "bundle.js",
	},
	// モジュールとして扱いたいファイルの拡張子を指定する
	// 例えば「import Foo from './foo'」という記述に対して"foo.ts"という名前のファイルをモジュールとして探す
	// デフォルトは['.js', '.json']
	resolve: {
		extensions: [".ts", ".js", ".tsx", ".jsx"],
	},
	devServer: {
		// webpack-dev-serverの公開フォルダ
		static: {
			directory: path.join(__dirname, "setting/dist"),
		},
	},
	// モジュールに適用するルールの設定（ここではローダーの設定を行う事が多い）
	module: {
		rules: [
			{
				// 拡張子が.tsで終わるファイルに対して、TypeScriptコンパイラを適用する
				test: /\.ts$/,
				loader: "ts-loader",
			},
			{
				test: /\.tsx$/,
				loader: "ts-loader",
			},
		],
	},
};
