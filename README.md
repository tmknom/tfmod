# tfmod

`tfmod` outputs the dependencies and dependents for Terraform.

## Description

`tfmod` helps you understand the relationships between your Terraform state directories and modules.
It shows which state directories use which modules and vice versa.

- `dependency`: Finds the modules used by the specified state directories.
- `dependent`: Finds the state directories that use the specified modules.

A "state directory" is where you run `terraform plan` or `terraform apply`.
It typically stores the `terraform.tfstate` file.

For example, if you modify a Terraform module, you can see which state directories are affected.

```shell
$ tfmod dependent --module=module/foo
state/prd state/stg
```

With `tfmod`, you can easily identify relationships and dependencies in your Terraform configurations,
making it easier to manage changes and understand their impacts.

## Installation

### Download binary

Download the latest compiled binaries and put it anywhere in your executable path.

- [GitHub Releases][releases]

### Pull container image

You can pull container image from GitHub Packages.

```shell
docker pull ghcr.io/tmknom/tfmod
```

## Usage

Assume you have the following directory structure for managing Terraform:

```md
├── state
│   ├── prd
│   └── stg
└── module
    ├── bar
    └── foo
```

You can run the subcommands as follows:

### Download Terraform modules

Download all the Terraform modules for the state directories under the current directory.

```shell
tfmod download
```

You can also specify the base directory for the state directories.

```shell
tfmod download --base state
```

> [!TIP]
>
> To analyze dependencies with `tfmod`, you need to download Terraform modules first.
> If you have already run `terraform get` or `terraform init`, this command is not needed.

### Explore dependencies

Specify a state directory to find the Terraform modules it depends on.

```shell
tfmod dependency --state=state/prd
```

You can specify multiple state directories separated by commas.

```shell
tfmod dependency --state=state/prd,state/stg
```

### Explore dependents

Specify a Terraform module to find the state directories that reference it.

```shell
tfmod dependent --module=module/foo
```

You can specify multiple Terraform modules separated by commas.

```shell
tfmod dependent --module=module/foo,module/bar
```

## Related projects

N/A

## Release notes

See [GitHub Releases][releases].

[releases]: https://github.com/tmknom/tfmod/releases/latest
