import json
from pathlib import Path
from typing import Optional


def read_save(source_file: Path) -> Optional[list[dict]]:
    save_dir = source_file.parent / ".lazycph"
    save_file = save_dir / f"{source_file.name}.json"
    if not save_file.exists():
        return None
    return json.loads(save_file.read_bytes())


def save_file(source_file: Path, data: list[dict]):
    base_dir = source_file.parent
    save_dir = base_dir / ".lazycph"
    if not save_dir.exists():
        save_dir.mkdir()
        save_dir.joinpath(".gitignore").write_text("*\n")
    save_file = save_dir / f"{source_file.name}.json"
    with save_file.open("w") as f:
        json.dump(data, f)
