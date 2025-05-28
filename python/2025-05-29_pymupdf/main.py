import json
import pymupdf  # type: ignore
from pymupdf import utils as pymupdfutils
from pathlib import Path


def get_text(pdf_path: Path) -> list[str]:
    doc = pymupdf.Document(pdf_path)

    texts: list[str] = []
    for page in doc:
        # NOTE: page.get_text("text") is dynamically assigned, so it's not recognized by type checkers.
        # https://github.com/pymupdf/PyMuPDF/issues/2883
        # text = page.get_text("text")  # type: ignore
        text = str(pymupdfutils.get_text(page, "text"))
        texts.append(text)
    doc.close()

    return texts


def highlight_text(pdf_path: Path, search_text: str) -> pymupdf.Document:
    highlighted_doc = pymupdf.Document(pdf_path)
    for page in highlighted_doc:
        quads = pymupdfutils.search_for(page, search_text, quads=True)
        page.add_highlight_annot(quads)

    return highlighted_doc


def save_document(doc: pymupdf.Document, output_path: Path) -> None:
    doc.save(output_path, garbage=4, deflate=True, clean=True)
    doc.close()


if __name__ == "__main__":
    output_dir = Path("output")
    output_dir.mkdir(exist_ok=True)

    pdf_url = "https://github.com/pymupdf/PyMuPDF/blob/c469893cdc335c56ec07c4bbbebc6e31ea25e01e/tests/resources/2201.00069.pdf?raw=true"
    pdf_path = Path("example.pdf")
    if not pdf_path.exists():
        import requests  # type: ignore

        print(f"Downloading PDF from {pdf_url}...")
        response = requests.get(pdf_url)
        with open(pdf_path, "wb") as f:
            f.write(response.content)

    # Get text from the PDF
    text = get_text(pdf_path)
    with open(output_dir / "texts.json", "w", encoding="utf-8") as f2:
        json.dump(text, f2, ensure_ascii=False, indent=4)

    # Highlight text in the PDF
    search_text = """Compact persistent emission was detected in the 1.51 GHz e-
MERLIN image at R.A. = 12ℎ15𝑚55𝑠.116, Dec. = −13◦01′14.′′48
at 86 𝜇Jy beam−1 by e-MERLIN.""".replace("-\n", "")
    highlighted_doc = highlight_text(pdf_path, search_text)

    output_path = output_dir / pdf_path.with_name(f"{pdf_path.stem}_highlighted.pdf")
    save_document(highlighted_doc, output_path)
    print(f"Highlighted PDF saved to {output_path}")
