with (import <nixpkgs> {config.allowUnfree = true;});
mkShell {
  shellHook = ''
    export PATH="$PATH:$(go env GOPATH)/bin"
    '';
  buildInputs = [
      go
      gnumake
  ];
}
