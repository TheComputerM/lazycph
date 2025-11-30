import argparse
from pathlib import Path

from lazycph import __version__
from lazycph.app import LazyCPH


def validate_target_path(path_str):
    """Validate that the target path exists."""
    path = Path(path_str).resolve()
    if not path.exists():
        raise argparse.ArgumentTypeError(f"Path '{path_str}' does not exist")
    return path


def parse_arguments():
    parser = argparse.ArgumentParser(
        prog="LazyCPH", description="Competitive Programming Helper"
    )
    parser.add_argument(
        "-v",
        "--version",
        action="version",
        version=f"%(prog)s {__version__}",
    )
    parser.add_argument(
        "target",
        nargs="?",
        default=".",
        type=validate_target_path,
        help="Target directory or file",
    )
    parser.add_argument(
        "-c",
        "--companion",
        action="store_true",
        default=False,
        help="Enable companion mode",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_arguments()
    assert isinstance(args.target, Path)

    base, selected = (
        (args.target.parent, args.target)
        if args.target.is_file()
        else (args.target, None)
    )

    app = LazyCPH(base, selected, args.companion)
    app.run()


if __name__ == "__main__":
    main()
