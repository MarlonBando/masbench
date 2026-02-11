#!/usr/bin/env python3
"""
Release script for masbench
This script automates the release process by:
1. Reading and validating the VERSION file
2. Extracting changelog from docs/source/changes.rst
3. Running pre-flight checks (git status, unpushed commits, tag existence)
4. Building cross-platform binaries
5. Creating git tag and GitHub release
"""

import os
import re
import subprocess
import sys
from pathlib import Path


class Colors:
    """ANSI color codes for terminal output"""
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    NC = '\033[0m'  # No Color


def print_error(msg):
    """Print error message in red"""
    print(f"{Colors.RED}✗ Error: {msg}{Colors.NC}", file=sys.stderr)


def print_success(msg):
    """Print success message in green"""
    print(f"{Colors.GREEN}✓ {msg}{Colors.NC}")


def print_info(msg):
    """Print info message in blue"""
    print(f"{Colors.BLUE}ℹ {msg}{Colors.NC}")


def print_warning(msg):
    """Print warning message in yellow"""
    print(f"{Colors.YELLOW}⚠ {msg}{Colors.NC}")


def run_command(cmd, capture_output=True, check=True):
    """Run shell command and return result"""
    try:
        result = subprocess.run(
            cmd,
            shell=True,
            capture_output=capture_output,
            text=True,
            check=check
        )
        return result
    except subprocess.CalledProcessError as e:
        print_error(f"Command failed: {cmd}")
        if e.stderr:
            print(e.stderr, file=sys.stderr)
        sys.exit(1)


def read_version():
    """Read and validate version from VERSION file"""
    version_file = Path("VERSION")
    
    if not version_file.exists():
        print_error("VERSION file not found")
        sys.exit(1)
    
    version = version_file.read_text().strip()
    
    # Validate semver format (X.Y.Z)
    if not re.match(r'^\d+\.\d+\.\d+$', version):
        print_error(f"Invalid version format '{version}'. Expected semver (e.g., 1.2.3)")
        sys.exit(1)
    
    print_success(f"Version: {version}")
    return version


def extract_changelog(version):
    """Extract changelog for the given version from changes.rst"""
    changes_file = Path("docs/source/changes.rst")
    
    if not changes_file.exists():
        print_error("docs/source/changes.rst not found")
        sys.exit(1)
    
    content = changes_file.read_text()
    
    # Look for version section (e.g., "Version 1.2.1" followed by dashes)
    version_pattern = rf'^Version {re.escape(version)}\s*\n-+\s*\n'
    match = re.search(version_pattern, content, re.MULTILINE)
    
    if not match:
        print_error(f"No changes found for version {version} in docs/source/changes.rst")
        print_info("Make sure the changelog has a section like:")
        print(f"  Version {version}")
        print(f"  {'-' * (8 + len(version))}")
        sys.exit(1)
    
    # Extract content from this version until the next "Version" heading or EOF
    start_pos = match.end()
    
    # Find next version section
    next_version_match = re.search(r'\nVersion \d+\.\d+\.\d+\s*\n-+', content[start_pos:])
    
    if next_version_match:
        end_pos = start_pos + next_version_match.start()
        changelog = content[start_pos:end_pos].strip()
    else:
        changelog = content[start_pos:].strip()
    
    if not changelog:
        print_error(f"Changelog for version {version} is empty")
        sys.exit(1)
    
    print_success(f"Extracted changelog for version {version}")
    return changelog


def rst_to_markdown(rst_text):
    """Convert basic reStructuredText to Markdown"""
    # Convert **text:** to **text:**
    md = rst_text
    
    # Convert * bullet points (already markdown compatible)
    # Convert ** subsections to ### headers
    md = re.sub(r'\*\*([^*]+):\*\*', r'**\1:**', md)
    
    # Convert ``code`` to `code`
    md = re.sub(r'``([^`]+)``', r'`\1`', md)
    
    return md


def check_git_status():
    """Check if working directory is clean"""
    result = run_command("git status --porcelain")
    
    if result.stdout.strip():
        print_error("Working directory is not clean. Commit or stash your changes first.")
        print("\nUncommitted changes:")
        print(result.stdout)
        sys.exit(1)
    
    print_success("Working directory is clean")


def check_unpushed_commits():
    """Check for unpushed commits"""
    # Get current branch
    result = run_command("git rev-parse --abbrev-ref HEAD")
    branch = result.stdout.strip()
    
    # Check if there are unpushed commits
    result = run_command(f"git log origin/{branch}..HEAD --oneline", check=False)
    
    if result.returncode == 0 and result.stdout.strip():
        unpushed_commits = result.stdout.strip().split('\n')
        count = len(unpushed_commits)
        
        print_warning(f"You have {count} unpushed commit(s) on branch '{branch}':")
        for commit in unpushed_commits:
            print(f"  {commit}")
        
        response = input(f"\nPush them first? [y/N]: ").strip().lower()
        if response == 'y':
            print_info(f"Pushing commits to origin/{branch}...")
            run_command(f"git push origin {branch}", capture_output=False)
            print_success("Commits pushed successfully")
        else:
            print_warning("Continuing without pushing commits...")
    else:
        print_success("No unpushed commits")


def check_tag_exists(version):
    """Check if git tag already exists"""
    tag = f"v{version}"
    result = run_command(f"git tag -l {tag}")
    
    if result.stdout.strip():
        print_error(f"Git tag '{tag}' already exists")
        print_info("Either:")
        print(f"  1. Update VERSION file to a new version")
        print(f"  2. Delete the existing tag: git tag -d {tag}")
        sys.exit(1)
    
    print_success(f"Tag 'v{version}' does not exist")


def check_gh_cli():
    """Check if gh CLI is installed and authenticated"""
    result = run_command("which gh", check=False)
    
    if result.returncode != 0:
        print_error("GitHub CLI (gh) is not installed")
        print_info("Install it from: https://cli.github.com/")
        sys.exit(1)
    
    result = run_command("gh auth status", check=False)
    
    if result.returncode != 0:
        print_error("GitHub CLI is not authenticated")
        print_info("Run: gh auth login")
        sys.exit(1)
    
    print_success("GitHub CLI is installed and authenticated")


def build_binaries(version):
    """Build cross-platform binaries and package them into zip files"""
    import zipfile
    import shutil
    
    print_info("Building cross-platform binaries...")
    
    # Create dist directory
    dist_dir = Path("dist")
    if dist_dir.exists():
        shutil.rmtree(dist_dir)
    dist_dir.mkdir(exist_ok=True)
    
    platforms = [
        ("linux", "amd64", ""),
        ("linux", "arm64", ""),
        ("windows", "amd64", ".exe"),
        ("windows", "arm64", ".exe"),
        ("darwin", "amd64", ""),
        ("darwin", "arm64", ""),
    ]
    
    zip_files = []
    
    for goos, goarch, ext in platforms:
        # Create platform-specific directory
        platform_dir = dist_dir / f"{goos}-{goarch}"
        platform_dir.mkdir(exist_ok=True)
        
        # Simple binary name (masbench or masbench.exe)
        binary_name = f"masbench{ext}"
        binary_path = platform_dir / binary_name
        
        print(f"  Building {goos}/{goarch}...")
        
        env = os.environ.copy()
        env['GOOS'] = goos
        env['GOARCH'] = goarch
        
        result = subprocess.run(
            ["go", "build", "-o", str(binary_path), "."],
            env=env,
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            print_error(f"Failed to build {goos}/{goarch}")
            print(result.stderr, file=sys.stderr)
            sys.exit(1)
        
        # Create zip file
        zip_name = f"masbench-v{version}-{goos}-{goarch}.zip"
        zip_path = dist_dir / zip_name
        
        print(f"  Creating {zip_name}...")
        
        with zipfile.ZipFile(zip_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
            zipf.write(binary_path, binary_name)
        
        # Clean up the platform directory
        shutil.rmtree(platform_dir)
        
        zip_files.append(zip_path)
    
    print_success(f"Built and packaged {len(zip_files)} binaries")
    return zip_files


def create_release(version, changelog, binaries):
    """Create git tag and GitHub release"""
    tag = f"v{version}"
    
    # Create git tag
    print_info(f"Creating git tag '{tag}'...")
    run_command(f"git tag {tag}")
    print_success(f"Tag '{tag}' created")
    
    # Push tag
    print_info(f"Pushing tag to origin...")
    run_command(f"git push origin {tag}", capture_output=False)
    print_success(f"Tag pushed to origin")
    
    # Convert changelog to markdown
    changelog_md = rst_to_markdown(changelog)
    
    # Create release notes
    release_notes = f"""## masbench v{version}

### Installation

Download the appropriate zip file for your platform below, extract it, and you're ready to go!

- **Linux AMD64**: `masbench-v{version}-linux-amd64.zip`
- **Linux ARM64**: `masbench-v{version}-linux-arm64.zip`
- **Windows AMD64**: `masbench-v{version}-windows-amd64.zip`
- **Windows ARM64**: `masbench-v{version}-windows-arm64.zip`
- **macOS Intel**: `masbench-v{version}-darwin-amd64.zip`
- **macOS Apple Silicon**: `masbench-v{version}-darwin-arm64.zip`

### Quick Start

1. Download the zip file for your platform
2. Extract it: `unzip masbench-v{version}-<platform>.zip`
3. Make it executable (Linux/macOS): `chmod +x masbench`
4. Move to your PATH: `sudo mv masbench /usr/local/bin/`
5. Run: `masbench --help`

### Changes

{changelog_md}
"""
    
    # Write release notes to temp file
    notes_file = Path("release_notes.md")
    notes_file.write_text(release_notes)
    
    # Create GitHub release
    print_info(f"Creating GitHub release...")
    
    cmd_parts = [
        "gh", "release", "create", tag,
        *[str(b) for b in binaries],
        "--title", f"masbench v{version}",
        "--notes-file", str(notes_file),
        "--latest"
    ]
    
    result = subprocess.run(cmd_parts, capture_output=True, text=True)
    
    if result.returncode != 0:
        print_error("Failed to create GitHub release")
        print(result.stderr, file=sys.stderr)
        # Clean up tag
        print_info("Cleaning up tag...")
        run_command(f"git tag -d {tag}", check=False)
        run_command(f"git push origin :refs/tags/{tag}", check=False)
        sys.exit(1)
    
    # Clean up temp file
    notes_file.unlink()
    
    print_success(f"GitHub release v{version} created successfully!")
    print_info(f"View it at: https://github.com/{get_repo_slug()}/releases/tag/{tag}")


def get_repo_slug():
    """Get repository slug (owner/repo) from git remote"""
    result = run_command("git config --get remote.origin.url")
    url = result.stdout.strip()
    
    # Extract owner/repo from URL
    match = re.search(r'github\.com[:/]([^/]+/[^/]+?)(\.git)?$', url)
    if match:
        return match.group(1)
    
    return "OWNER/REPO"


def main():
    """Main function"""
    print(f"\n{Colors.BLUE}{'=' * 60}{Colors.NC}")
    print(f"{Colors.BLUE}  masbench Release Script{Colors.NC}")
    print(f"{Colors.BLUE}{'=' * 60}{Colors.NC}\n")
    
    # Step 1: Read and validate version
    print(f"{Colors.BLUE}[1/5] Reading VERSION file...{Colors.NC}")
    version = read_version()
    print()
    
    # Step 2: Extract changelog
    print(f"{Colors.BLUE}[2/5] Extracting changelog...{Colors.NC}")
    changelog = extract_changelog(version)
    print()
    
    # Step 3: Pre-flight checks
    print(f"{Colors.BLUE}[3/5] Running pre-flight checks...{Colors.NC}")
    check_git_status()
    check_unpushed_commits()
    check_tag_exists(version)
    check_gh_cli()
    print()
    
    # Step 4: Build binaries
    print(f"{Colors.BLUE}[4/5] Building binaries...{Colors.NC}")
    binaries = build_binaries(version)
    print()
    
    # Step 5: Create release
    print(f"{Colors.BLUE}[5/5] Creating release...{Colors.NC}")
    create_release(version, changelog, binaries)
    print()
    
    print(f"{Colors.GREEN}{'=' * 60}{Colors.NC}")
    print(f"{Colors.GREEN}  Release v{version} completed successfully!{Colors.NC}")
    print(f"{Colors.GREEN}{'=' * 60}{Colors.NC}\n")


if __name__ == "__main__":
    main()
