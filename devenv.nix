{pkgs, ...}: {
  packages = [pkgs.git];

  languages.rust.enable = true;

  pre-commit.hooks = {
    rustfmt.enable = true;
    clippy.enable = true;
  };
}
