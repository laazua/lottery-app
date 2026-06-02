### lottery-app

---
### 说明
1. 大乐透随机号码

---
##### 依赖
1. jdk1.8
2. gioui.org
3. commandlinetools
4. 环境配置
```bash
# 1. 安装go环境
# 2. 安装jdk1.8环境
# 3. 安装commandlinetools
wget https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip
unzip commandlinetools-linux-*.zip
mkdir -p ~/Android/Sdk
mv cmdline-tools ~/Android/Sdk/

# 设置环境变量
echo 'export ANDROID_HOME=$HOME/Android/Sdk' >> ~/.bashrc
echo 'export PATH=$PATH:$ANDROID_HOME/cmdline-tools/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 SDK 和 NDK (接受协议)
sdkmanager "platforms;android-31" "ndk-bundle" "build-tools;31.0.0"

# 4. 安装gio工具
go install gioui.org/cmd/gogio@latest
# 确保 $GOPATH/bin 在 PATH 中
export PATH=$PATH:$HOME/go/bin
```

---
##### 打包
```bash
gogio -target android -appid com.lottery.app -icon lottery.png -o lottery.apk .
```

---
##### [表准库移动开发](https://golang.org/x/mobile)
