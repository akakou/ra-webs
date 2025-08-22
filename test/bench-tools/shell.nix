let pkgs = import <nixpkgs> {};

in pkgs.mkShell rec {
  name = "ra-webs-bench";

  buildInputs = with pkgs; [
    python3
    python3Packages.selenium
    chromedriver
    chromium
  ];
}