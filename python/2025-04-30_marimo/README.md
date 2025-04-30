# 2025-04-20_marimo

Official page: [marimo](https://docs.marimo.io/)

> marimo is a reactive Python notebook

Jupyter notebookの代替として挙げられていたPythonライブラリ。
`.py`で書ける上にSliderやAccordionなどのコンポーネントが含まれているのが特徴。
WASMを使ってWebアプリとしてデプロイもできるらしい。

[Online playground](https://links.marimo.app/tutorial-intro)からすぐ試すことができる。

[Jupyter code cells](https://code.visualstudio.com/docs/python/jupyter-support-py#_jupyter-code-cells)のようにJupyter NotebookをPythonで書けるようにするものかと思ったが、Web上に編集用サーバーを建ててコード/ドキュメントを書く体験を提供するライブラリだった。
Web上で記述してセーブした後のPythonファイルは以下のようにmarimoに依存した形になる。

```python
app = marimo.App()

@app.cell
def _():
    import marimo as mo

    mo.md("# Welcome to marimo! 🌊🍃")
    return (mo,)

if __name__ == "__main__":
    app.run()
```

## Links

- [次世代notebook『marimo』入門（#13） - てっく・ざ・ぶろぐ！](https://alvinvin.hatenablog.jp/entry/13)
