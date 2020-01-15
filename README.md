<img alt="Terrier Logo" src="https://user-images.githubusercontent.com/449385/72431923-d0219680-378d-11ea-82e4-959da8ebbc9e.png" width="256" />

# Terrier
Terrier is a Image and Container analysis tool that can be used to scan OCI images and Containers to identify and verify the presence of specific files according to their hashes. A detailed writeup of Terrier can be found on the Heroku blog, https://blog.heroku.com/terrier-open-source-identifying-analyzing-containers.

# Installation
#### Binaries
For installation instructions from binaries please visit the [Releases Page](https://github.com/heroku/terrier/releases).

#### Via Go

```
$ go get github.com/heroku/terrier
```

## Building from source

Via go

```
$ go build
```

or 

```
$ make all
```

## Usage

```console
$ ./terrier -h
Usage of ./terrier:
  -cfg string
        Load config from provided yaml file (default "cfg.yml")
```


An OCI TAR of the image to be scanned is required, this is provided to Terrier via the "Image" value in the cfg.yml. 

The following Docker command can be used to convert a Docker image to a TAR that can be scanned by Terrier.

```console
# docker save imageid -o image.tar

```

```console
$ ./terrier 
[+] Loading config: cfg.yml
[+] Analysing Image
[+] Docker Image Source:  image.tar
[*] Inspecting Layer:  05c3c2c60920f68b506d3c66e0f6148b81a8b0831388c2d61be5ef02190bcd1f
[!] All components were identified and verified: (493/493)
```

### Example YML config
Terrier parses YAML, below is an example config.

```
#THIS IS AN EXAMPLE CONFIG, MODIFY TO YOUR NEEDS

mode: image
image: image.tar
# mode: container
# path: merged
# verbose: true
# veryverbose: true

files:
  - name: '/usr/bin/curl'
    hashes:
      - hash: '2353cbb7b47d0782ba8cdd9c7438b053c982eaaea6fbef8620c31a58d1e276e8'
      - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2aaa'
      - hash: '9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96'
      - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521faa'
  - name: '/usr/bin/go'
    hashes:
      - hash: '2353cbb7b47d0782ba8cdd9c7438b053c982eaaea6fbef8620c31a58d1e276e8'

#UNCOMMENT TO ANALYZE HASHES
#     hashes:
#       - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521faa'
#       - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2aa'
#       - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41aa'

```

# What does Terrier do?
Terrier is a CLI tool that allows you to:
- Scan an OCI image for the presence of one or more files that match one or more provided SHA256 hashes
- Scan a running Container for the presence of one or more files that match one or more provided SHA256 hashes

# What is Terrier useful for?
### Scenario 1
Terrier can be used to verify if a specific OCI image is making use of a specific binary, which useful in a supply chain verification scenario. For example, we may want to check that a specific Docker image is making use of a specific version or versions of cURL. In this case, Terrier is supplied with the SHA256 hashes of the binaries that are trusted.

An example YAML file for this scenario might look like this:
```
mode: image
# verbose: true
# veryverbose: true
image: golang1131.tar

files:
  - name: '/usr/local/bin/analysis.sh'
    hashes:
       - hash: '9adc0bf7362bb66b98005aebec36691a62c80d54755e361788c776367d11b105'
  - name: '/usr/bin/curl'
    hashes:
       - hash: '23afbfab4f35ac90d9841a6e05f0d1487b6e0c3a914ea8dab3676c6dde612495'
  - name: '/usr/local/bin/staticcheck'
    hashes:
       - hash: '73f89162bacda8dd2354021dc56dc2f3dba136e873e372312843cd895dde24a2'
```


### Scenario 2
Terrier can be used to verify the presence of a particular file or files in a OCI image according to a set of provided hashes. This can be useful to check if an OCI image contains a malicious file or a file that is required to be identified.

An example YAML file for this scenario might look like this:
```
mode: image
# verbose: true
# veryverbose: true
image: alpinetest.tar
hashes:
  - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f'
  - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2'
  - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41'
  - hash: '9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96'
```


### Scenario 3
Terrier can be used to verify the components of Containers at runtime by analysing the contents of /var/lib/docker/overlay2/.../merged 
An example YAML file for this scenario might look like this:
```
mode: container
verbose: true
# veryverbose: true
# image: latestgo13.tar
path: merged

files:
  - name: '/usr/local/bin/analysis.sh'
    hashes:
       - hash: '9adc0bf7362bb66b98005aebec36691a62c80d54755e361788c776367d11b105'
  - name: '/usr/local/go/bin/go'
    hashes:
       - hash: '23afbfab4f35ac90d9841a6e05f0d1487b6e0c3a914ea8dab3676c6dde612495'
  - name: '/usr/local/bin/staticcheck'
    hashes:
       - hash: '73f89162bacda8dd2354021dc56dc2f3dba136e873e372312843cd895dde24a2'
  - name: '/usr/local/bin/gosec'
    hashes:
       - hash: 'e7cb8304e032ccde8e342a7f85ba0ba5cb0b8383a09a77ca282793ad7e9f8c1f'
  - name: '/usr/local/bin/errcheck'
    hashes:
       - hash: '41f725d7a872cad4ce1f403938937822572e0a38a51e8a1b29707f5884a2f0d7'
  - name: '/var/lib/dpkg/info/apt.postrm'
    hashes:
       - hash: '6a8f9af3abcfb8c6e35887d11d41a83782b50f5766d42bd1e32a38781cba0b1c'
```



# Usage
### Example 1
Terrier is a CLI and makes use of YAML. An example YAML config:
```
mode: image
# verbose: true
# veryverbose: true
image: alpinetest.tar
files:
  - name: '/usr/local/go/bin/go'
    hashes:
      - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f'
      - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2aaa'
      - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41aaa'
      - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521faa'
  - name: '/usr/bin/delpart'
    hashes:
      - hash: '9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96aaa'
  - name: '/usr/bin/stdbuf'
    hashes:
      - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521faa'
      - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2aa'
      - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41aa'
```

In the example below, Terrier has being instructed via the YAML above to verify multiple files.
```
$./terrier 
[+] Loading config: cfg.yml
[+] Analysing Image
[+] Docker Image Source:  alpinetest.tar
[*] Inspecting Layer:  05c3c2c60920f68b506d3c66e0f6148b81a8b0831388c2d61be5ef02190bcd1f
[*] Inspecting Layer:  09c25a178d8a6f8b984f3e72ca5ec966215b24a700ed135dc062ad925aa5eb23
[*] Inspecting Layer:  36351e8e1da92268d40245cfbcd499a1173eeacc23be428386c8fc0a16f0b10a
[*] Inspecting Layer:  7224ca1e886eeb7e63a9e978b1a811ed52f4a53ccb65f7c510fa04a0d1103fdf
[*] Inspecting Layer:  7a2e464d80c7a1d89dab4321145491fb94865099c59975cfc840c2b8e7065014
[*] Inspecting Layer:  88a583fe02f250344f89242f88309c666671042b032411630de870a111bea971
[*] Inspecting Layer:  8db14b6fdd2cf8b4c122824531a4d85e07f1fecd6f7f43eab7f2d0a90d8c4bf2
[*] Inspecting Layer:  9196e3376d1ed69a647e728a444662c10ed21feed4ef7aaca0d10f452240a09a
[*] Inspecting Layer:  92db9b9e59a64cdf486203189d02acff79c3360788b62214a49d2263874ee811
[*] Inspecting Layer:  bc4bb4a45da628724c9f93400a9149b2dd8a5d437272cb4e572cfaec64512d98
[*] Inspecting Layer:  be7d600e4e8ed3000e342ef6482211350069d935a14aeff4d9fc3289e1426ed3
[*] Inspecting Layer:  c4cec85dfa44f0a8856064922cff1c39b872b506dd002e33664d11a80f75a149
[*] Inspecting Layer:  c998d6f023b7b9e3c186af19bcd1c2574f0d01b943077281ac5bd32e02dc57a5
[!] All components were identified and verified: (493/493)

```
 Terrier sets its return code depending on the result of the tests, in the case of the test above, the return code will be "0" which indicates a successful test as 1 instance of each provided component was identified and verified.
 
### Example 2
Terrier is instructed to identify any files in the provided image that match the provided SHA256 hashes.
YAML file cfg.yml
```
mode: image
# verbose: true
# veryverbose: true
image: 1070caa1a8d89440829fd35d9356143a9d6185fe7f7a015b992ec1d8aa81c78a.tar
hashes:
  - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f'
  - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2'
  - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41'
  - hash: '9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96'
```
Running Terrier.
```
./terrier 
[+] Loading config: cfg.yml
[+] Docker Image Source:  golang.tar
[*] Inspecting Layer:  1070caa1a8d89440829fd35d9356143a9d6185fe7f7a015b992ec1d8aa81c78a
[*] Inspecting Layer:  414833cdb33683ab8607565da5f40d3dc3f721e9a59e14e373fce206580ed40d
[*] Inspecting Layer:  6bd93c6873c822f793f770fdf3973d8a02254a5a0d60d67827480797f76858aa
[*] Inspecting Layer:  c40c240ae37a2d2982ebcc3a58e67bf07aeaebe0796b5c5687045083ac6295ed
[*] Inspecting Layer:  d2850df0b6795c00bdce32eb9c1ad9afc0640c2b9a3e53ec5437fc5539b1d71a
[*] Inspecting Layer:  f0c2fe7dbe3336c8ba06258935c8dae37dbecd404d2d9cd74c3587391a11b1af
        [!] Found file 'f0c2fe7dbe3336c8ba06258935c8dae37dbecd404d2d9cd74c3587391a11b1af/usr/bin/curl' with hash: 9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96
[*] Inspecting Layer:  f2d913644763b53196cfd2597f21b9739535ef9d5bf9250b9fa21ed223fc29e3
echo $?
1
```

### Example 3
Terrier is instructed to analyze and verify the contents of the container's merged contents located at "merged" where merged is possibly located at ```/var/lib/docker/overlay2/..../merged ```.
An example YAML file for this scenario might look like this:
```
mode: container
verbose: true
# veryverbose: true
# image: latestgo13.tar
path: merged

files:
  - name: '/usr/local/bin/analysis.sh'
    hashes:
       - hash: '9adc0bf7362bb66b98005aebec36691a62c80d54755e361788c776367d11b105'
  - name: '/usr/local/go/bin/go'
    hashes:
       - hash: '23afbfab4f35ac90d9841a6e05f0d1487b6e0c3a914ea8dab3676c6dde612495'
  - name: '/usr/local/bin/staticcheck'
    hashes:
       - hash: '73f89162bacda8dd2354021dc56dc2f3dba136e873e372312843cd895dde24a2'
  - name: '/usr/local/bin/gosec'
    hashes:
       - hash: 'e7cb8304e032ccde8e342a7f85ba0ba5cb0b8383a09a77ca282793ad7e9f8c1f'
  - name: '/usr/local/bin/errcheck'
    hashes:
       - hash: '41f725d7a872cad4ce1f403938937822572e0a38a51e8a1b29707f5884a2f0d7'
  - name: '/var/lib/dpkg/info/apt.postrm'
    hashes:
       - hash: '6a8f9af3abcfb8c6e35887d11d41a83782b50f5766d42bd1e32a38781cba0b1c'
```
Running Terrier to analyse the running Container.

```
[+] Loading config: cfg.yml
[+] Analysing Container
[!] Found matching instance of '/usr/local/bin/analysis.sh' at: merged/usr/local/bin/analysis.sh with hash:9adc0bf7362bb66b98005aebec36691a62c80d54755e361788c776367d11b105
[!] Found matching instance of '/usr/local/bin/errcheck' at: merged/usr/local/bin/errcheck with hash:41f725d7a872cad4ce1f403938937822572e0a38a51e8a1b29707f5884a2f0d7
[!] Found matching instance of '/usr/local/bin/gosec' at: merged/usr/local/bin/gosec with hash:e7cb8304e032ccde8e342a7f85ba0ba5cb0b8383a09a77ca282793ad7e9f8c1f
[!] Found matching instance of '/usr/local/bin/staticcheck' at: merged/usr/local/bin/staticcheck with hash:73f89162bacda8dd2354021dc56dc2f3dba136e873e372312843cd895dde24a2
[!] Found matching instance of '/usr/local/go/bin/go' at: merged/usr/local/go/bin/go with hash:23afbfab4f35ac90d9841a6e05f0d1487b6e0c3a914ea8dab3676c6dde612495
[!] Found matching instance of '/var/lib/dpkg/info/apt.postrm' at: merged/var/lib/dpkg/info/apt.postrm with hash:6a8f9af3abcfb8c6e35887d11d41a83782b50f5766d42bd1e32a38781cba0b1c
[!] All components were identified and verified: (6/6)
```
 
 # Integrating with CI
 Terrier has been designed to assist in the prevention of supply chain attacks. To utilise Terrier with CI's such as Github actions or CircleCI, the following example configurations might be useful.
  
 ## CircleCI Example
 config.yml
 ```
version: 2
jobs:
 build:
   machine: true
   steps:
     - checkout
     - run:
        name: Build Docker Image
        command: |
              docker build -t builditall .
     - run:
        name: Save Docker Image Locally
        command: |
              docker save builditall -o builditall.tar
     - run:
        name: Verify Docker Image Binaries
        command: |
              ./terrier_linux_amd64
 ```
   Terrier cfg.yml
   ```
   mode:image
   image: builditall.tar
files:
  - name: '/bin/wget'
    hashes:
      - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f'
      - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2a'
      - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41a'
  - name: '/sbin/sulogin'
    hashes:
      - hash: '9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96aaa'
   ```
   
   
   
  ## Github Actions Example
  go.yml
  ```
  name: Go
on: [push]
jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:

    - name: Get Code
      uses: actions/checkout@master
    - name: Build Docker Image
      run: |
        docker build -t builditall .
    - name: Save Docker Image Locally
      run: |
        docker save builditall -o builditall.tar
    - name: Verify Docker Image Binaries
      run: |
        ./terrier_linux_amd64
  ```

Terrier cfg.yml
   ```
mode: image
image: builditall.tar
files:
  - name: '/bin/wget'
    hashes:
      - hash: '8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f'
      - hash: '22e88c7d6da9b73fbb515ed6a8f6d133c680527a799e3069ca7ce346d90649b2a'
      - hash: '60a2c86db4523e5d3eb41a247b4e7042a21d5c9d483d59053159d9ed50c8aa41a'
  - name: '/bin/sbin/sulogin'
    hashes:
      - hash: '9a43cb726fef31f272333b236ff1fde4beab363af54d0bc99c304450065d9c96aaa'
   ```
 # Converting SHASUM 256 Hashes to a Terrier Config File
 Sometimes the source of SHA256 hashes is produced from other tools in the following format:

 ```
 6a8f9af3abcfb8c6e35887d11d41a83782b50f5766d42bd1e32a38781cba0b1c  ./var/lib/dpkg/info/apt.postrm
6374f7996297a6933c9ccae7eecc506a14c85112bf1984c12da1f975dab573b2  ./var/lib/dpkg/info/mawk.postinst
fd72e78277680d02dcdb5d898fc9e3fed00bf011ccf31deee0f9e5f4cf299055  ./var/lib/dpkg/info/lsb-base.preinst
fd72e78277680d02dcdb5d898fc9e3fed00bf011ccf31deee0f9e5f4cf299055  ./var/lib/dpkg/info/lsb-base.postrm
8a278d8f860ef64ae49a2d3099b698c79dd5184db154fdeaea1bc7544c2135df  ./var/lib/dpkg/info/debconf.postrm
1e6edefb6be6eb6fe8dd60ece5544938197b2d1d38a2d4957c069661bc2591cd  ./var/lib/dpkg/info/base-files.prerm
198c13dfc6e7ae170b48bb5b997793f5b25541f6e998edaec6e9812bc002915f  ./var/lib/dpkg/info/passwd.postinst
 ```

 The format above contains the data we need for Terrier but is in the wrong format. We have included a script called ```convertSHA.sh``` which can be used to convert a file with the file paths and hash values as seen above into a valid Terrier config file. 

 This can be seen in the following example:

 ```
 # cat hashes-SHA256.txt
 6a8f9af3abcfb8c6e35887d11d41a83782b50f5766d42bd1e32a38781cba0b1c  ./var/lib/dpkg/info/apt.postrm
6374f7996297a6933c9ccae7eecc506a14c85112bf1984c12da1f975dab573b2  ./var/lib/dpkg/info/mawk.postinst
fd72e78277680d02dcdb5d898fc9e3fed00bf011ccf31deee0f9e5f4cf299055  ./var/lib/dpkg/info/lsb-base.preinst
fd72e78277680d02dcdb5d898fc9e3fed00bf011ccf31deee0f9e5f4cf299055  ./var/lib/dpkg/info/lsb-base.postrm
8a278d8f860ef64ae49a2d3099b698c79dd5184db154fdeaea1bc7544c2135df  ./var/lib/dpkg/info/debconf.postrm
1e6edefb6be6eb6fe8dd60ece5544938197b2d1d38a2d4957c069661bc2591cd  ./var/lib/dpkg/info/base-files.prerm
198c13dfc6e7ae170b48bb5b997793f5b25541f6e998edaec6e9812bc002915f  ./var/lib/dpkg/info/passwd.postinst

# ./convertSHA.sh hashes-SHA256.txt output.yml
Converting hashes-SHA256.txt to Terrier YML: output.yml

# cat output.yml
mode: image
#mode: container
image: image.tar
#path: path/to/container/merged
#verbose: true
#veryverbose: true
files:
  - name: '/var/lib/dpkg/info/apt.postrm'
    hashes:
       - hash: '6a8f9af3abcfb8c6e35887d11d41a83782b50f5766d42bd1e32a38781cba0b1c'
  - name: '/var/lib/dpkg/info/mawk.postinst'
    hashes:
       - hash: '6374f7996297a6933c9ccae7eecc506a14c85112bf1984c12da1f975dab573b2'

 ```

