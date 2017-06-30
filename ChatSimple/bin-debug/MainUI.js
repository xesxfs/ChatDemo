var __reflect = (this && this.__reflect) || function (p, c, t) {
    p.__class__ = c, t ? t.push(c) : t = [c], p.__types__ = p.__types__ ? t.concat(p.__types__) : t;
};
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
// TypeScript file
var MainUI = (function (_super) {
    __extends(MainUI, _super);
    function MainUI() {
        var _this = _super.call(this) || this;
        _this.LoginHall = {
            msg: "Hello Lucy!!"
        };
        _this.Croom = {
            rid: 0
        };
        _this.Say = {
            msg: ""
        };
        _this.skinName = "MainSkin";
        return _this;
    }
    MainUI.prototype.childrenCreated = function () {
        this.createSocket();
        this.addEventListener("touchTap", this.onTouch, this);
    };
    MainUI.prototype.onTouch = function (e) {
        switch (e.target) {
            case this.loginBtn:
                var b = new egret.ByteArray();
                b.endian = egret.Endian.LITTLE_ENDIAN;
                b.writeUnsignedInt(1);
                b.writeUTFBytes(JSON.stringify(this.LoginHall));
                this.send(b);
                break;
            case this.croomBtn:
                var b = new egret.ByteArray();
                b.endian = egret.Endian.LITTLE_ENDIAN;
                b.writeUnsignedInt(2);
                // b.writeUTFBytes(JSON.stringify(this.LoginHall));
                this.send(b);
                break;
            case this.joinBtn:
                var b = new egret.ByteArray();
                b.endian = egret.Endian.LITTLE_ENDIAN;
                b.writeUnsignedInt(3);
                this.Croom.rid = parseInt(this.roomEdit.text);
                b.writeUTFBytes(JSON.stringify(this.Croom));
                this.send(b);
                break;
            case this.rinfoBtn:
                break;
            case this.sayBtn:
                var b = new egret.ByteArray();
                b.endian = egret.Endian.LITTLE_ENDIAN;
                b.writeUnsignedInt(5);
                this.Say.msg = this.sayEdit.text;
                b.writeUTFBytes(JSON.stringify(this.Say));
                this.send(b);
                break;
        }
    };
    MainUI.prototype.createSocket = function () {
        this.socket = new egret.WebSocket();
        this.socket.type = egret.WebSocket.TYPE_BINARY;
        this.socket.addEventListener(egret.Event.CONNECT, this.onConnect, this);
        this.socket.addEventListener(egret.Event.CLOSE, this.onClose, this);
        this.socket.addEventListener(egret.IOErrorEvent.IO_ERROR, this.onError, this);
        this.socket.addEventListener(egret.ProgressEvent.SOCKET_DATA, this.onRecieve, this);
        this.socket.connectByUrl("ws://127.0.0.1:8080");
    };
    MainUI.prototype.onClose = function (e) {
        console.log("close");
    };
    //连接错误
    MainUI.prototype.onError = function (e) {
    };
    //接收数据
    MainUI.prototype.onRecieve = function (e) {
        // console.log("--------------------------------------")
        var b = new egret.ByteArray();
        // b.endian = egret.Endian.LITTLE_ENDIAN;
        this.socket.readBytes(b);
        console.log(b.length);
        // var echo=this.socket.readUTF();
        var cmd = b.readUnsignedInt();
        var data = b.readUTFBytes(b.length - 4);
        console.log(cmd, data);
        // console.log(echo);
        // this.socket.close();
    };
    //连接成功
    MainUI.prototype.onConnect = function (e) {
        egret.log(this.name + " connect success");
    };
    MainUI.prototype.send = function (data) {
        this.socket.writeBytes(data);
        this.socket.flush();
    };
    return MainUI;
}(eui.Component));
__reflect(MainUI.prototype, "MainUI");
//# sourceMappingURL=MainUI.js.map