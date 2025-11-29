import argparse
from pathlib import Path

from lazycph.app import LazyCPH


def validate_target_path(path_str):
    """Validate that the target path exists."""
    path = Path(path_str)
    if not path.exists():
        raise argparse.ArgumentTypeError(f"Path '{path_str}' does not exist")
    return path


def parse_arguments():
    parser = argparse.ArgumentParser(
        prog="LazyCPH", description="Competitive Programming Helper"
    )
    parser.add_argument(
        "target",
        nargs="?",
        default=".",
        type=validate_target_path,
        help="Target directory or file",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_arguments()
    assert isinstance(args.target, Path)

    base = args.target if args.target.is_dir() else Path.cwd()
    selected = args.target if args.target.is_file() else None
    app = LazyCPH(base, selected)
    app.run()


if __name__ == "__main__":
    main()
