# Just build and pack your go application

The package is intended to provide an easy way to build and pack your go application for popular platforms.

There are other solutions like [goreleaser](https://goreleaser.com/) that are more complex and provide more features. 

This package is less flexible and provides fewer features, but is perfect if you need to build and pack your
go application for popular platforms with a single command which is a good choice for open-source projects.

Thanks to [xgo](https://github.com/crazy-max/xgo) it's fully os-independent and allows you to **cross-build** and pack 
an application for any supported platform and from any platform. **The only dependency is docker.**

## Supported targets

- MacOS: **amd64/arm64**
  - binary (+signing)
  - app bundle with icon (+signing)
  - dmg image with app bundle
- Windows: **amd64**
  - binary (with embed manifest + icon)
- Linux: **amd64/arm64**
  - binary
  - deb package (with desktop entry and icon)

## Installation

```bash
go install github.com/cardinalby/xgo-pack
```

## Create config file

The first step is creating `xgo-pack-config.yaml` file in the root of your project. To create a minimal working config file
you can use the following command:

```bash
xgo-pack init
```

## Start building and packing

```bash
xgo-pack build
```

It will use `xgo-pack-config.yaml` file in the current directory to build and pack your application.

## Check out the example

Example tray icon app that is built for all supported platforms is available in the [example](./example) directory.

### Config details

The config file is a yaml file. To see all available options use any modern IDE that supports json schema [available at
the repo](./config_schema/config.schema.v1.json), also 
you can refer to `Config` struct definition in [pkg/pipeline/config/cfgtypes/config.go](./pkg/pipeline/config/cfgtypes/config.go).

#### Config structure

Config has `common` options on different levels:
- `targets.common` is used as defaults for all os-es targets
- `targets.<os>.common` is used as for all architectures of the specified os and is based on `targets.common`

To enable building a particular artifact (target) you should set `targets.<os>.<arch>.build_<artifact>` to `true`:
- `targets.<os>.<arch>.build_bin` - build a binary
- `targets.macos.<arch>.build_bundle` - build an app bundle for MacOS targets
- `targets.macos.<arch>.build_dmg` - build a dmg image with app bundle for MacOS targets
- `targets.linux.<arch>.build_deb` - build deb package for Linux targets

#### Presets

You can use presets as template configs. 

```yaml
presets:
  - my-cfg-template-1.yaml
  - my-cfg-template-2.yaml
```

Templates are applied in the order they are listed. The last template has the highest priority.

**Built-in presets:**
- [`xgo-pack:gui`](./pkg/pipeline/config/presets/builtin/gui.go) - for GUI applications
- [`xgo-pack:tray_icon`](./pkg/pipeline/config/presets/builtin/tray_icon.go) - for Tray icon GUI applications

#### Icon
Icon is converted to the target formats using [ImageMagick](https://github.com/dooman87/imagemagick-docker). 
You can use src icon of any supported format. For psd files add `[0]` at the end to merge all layers

#### MacOS signing
[rcodesign](https://github.com/indygreg/apple-platform-rs) tool in docker is used for signing MacOS binaries and 
app bundles. By default, it signs them with self-signed certificate. 
To set up proper signing, add `rcodesign.toml` 
[config file](https://gregoryszorc.com/docs/apple-codesign/0.27.0/apple_codesign_rcodesign_config_files.html) 
to the root of your project.  

