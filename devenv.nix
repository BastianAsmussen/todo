{pkgs, ...}: {
  packages = with pkgs; [
    git
    go
  ];
}
