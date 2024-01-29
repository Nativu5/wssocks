# Instructions for Sangfor VPN

> **注意：** 本工作与深信服无关，且不以任何方式修改/逆向其产品，而是对其产品的合理使用。

为什么需要这个功能：

- 不想安装深信服的客户端
- 使用简单的服务（例如 SSH）

使用前提：

- 需要校内有一台可用的服务器

## 快速入门

1. 设置wssock服务器

    启动服务器时，应添加 `--twf` 选项启用对 Sangfor VPN 的支持。

    ```bash
    wssocks server --addr :1088 --twf
    ```

2. 登录 Sangfor VPN 以获取 TWF ID

    在 Sangfor Web VPN 中，wssock 服务器将被代理。例如：

    - wssock服务器地址：`example.com:1088`
    - 桑福Web VPN地址：`vpn.uestc.edu.cn`
    
    wssock 服务器将被代理到 `example-com-1088-p.vpn.uestc.edu.cn:8118`。

    > 注意：可能需要调整端口号 8118，因为不同机构的设置可能不同。

    访问被代理的地址时，您将被重定向到 Sangfor VPN 的登录页面。登录后，将显示 TWF ID。

3. 设置 wssock 客户端

    启动客户端时，添加 `--twfid` 选项启用对 Sangfor VPN 的支持。

    ```bash
    wssocks client --remote ws://example-com-1088-p.vpn.uestc.edu.cn:8118 --twfid YOUR_TWF_ID
    ```

## 高级用法

该功能不影响 wssock 的其他功能，请按照 [README.md](README.md) 中的说明使用其他功能。

## 参考资料

本分支的灵感来源于 [wssocks-plugin-ustb](https://github.com/genshen/wssocks-plugin-ustb)。

使用的是相同的原理，但使用 Sangfor VPN 的情况与 USTB 不同，因此实现方式有所差别。
