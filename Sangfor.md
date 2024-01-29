# Instructions for Sangfor VPN

[中文版](Sangfor_zh.md)

> **Note:** This work has no relationship with Sangfor Technologies Inc, and does not modify/reverse engineering Sangfor's products in any way, but is a fair use of its products.

## Quick Start

1. Set up the wssock server

    When you launch the server, you should add the `--twf` option to enable Sangfor VPN support.

    ```bash
    wssocks server --addr :1088 --twf
    ```

2. Login Sangfor VPN to get TWF ID

    In Sangfor Web VPN, the wssock server will be proxied. For example: 

    - wssock server address: `example.com:1088`
    - Sangfor Web VPN address: `vpn.uestc.edu.cn`

    Then the wssock server will be proxied to `example-com-1088-p.vpn.uestc.edu.cn:8118`.

    > **Note:** You may adjust the port number 8118 as it may be different in your institution.

    When visit the proxied address, you will be redirected to the login page of Sangfor VPN. After login, the TWF ID will be displayed. 

3. Set up the wssock client

    When you launch the client, you should add the `--twfid` option to enable Sangfor VPN support.

    ```bash
    wssocks client --remote ws://example-com-1088-p.vpn.uestc.edu.cn:8118 --twfid YOUR_TWF_ID
    ```

## Advanced Usage

Sangfor VPN support does not affect other features of wssocks. You can use other features as usual.

Regarding other settings, you can find detailed instructions in [README.md](README.md).

## References

This fork is inspired by [wssocks-plugin-ustb](https://github.com/genshen/wssocks-plugin-ustb). 

The same principle is used, but the case of Sangfor VPN is different from that of USTB. Therefore, the implementation is different.
