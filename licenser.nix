{ pkgs ? import <nixpkgs> {} }:

let
  name = "licenser";
  version = "0.7.0";
in
with pkgs;

buildGoModule {
  pname = name;
  version = "v${version}";
  src = ./.;
  vendorHash = "sha256-LBVVhg69VdQVsVARCkwooe6N6DacgViIW/iQWHCya4k=";
  ldFlags = "-w -s";
  CGO_ENABLED = "0";
  doCheck = false;
  meta = with lib; {
    description = "Licenser: Verify and prepend licenses to your GitHub repositories";
    homepage = "https://github.com/liamawhite/licenser";
    license = licenses.asl20;
    maintainers = [ maintainers.liamawhite ];
  };
}

