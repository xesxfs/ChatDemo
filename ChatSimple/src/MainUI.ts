// TypeScript file
class MainUI extends eui.Component{
    constructor(){
        super();
        this.skinName="MainSkin";    
    }

    private loginBtn:eui.Button;
    private croomBtn:eui.Button;
    private joinBtn:eui.Button;
    private rinfoBtn:eui.Button;
    private roomEdit:eui.EditableText;

    protected childrenCreated(){
        this.createSocket();        
        this.addEventListener("touchTap",this.onTouch,this);
    }

    public  LoginHall={
        msg:"Hello Lucy!!"
    }

    public Croom={
        rid:0
    }


    

    private onTouch(e:egret.TouchEvent){
        switch(e.target){
            case this.loginBtn:
                    var b: egret.ByteArray = new egret.ByteArray();
                    b.endian = egret.Endian.LITTLE_ENDIAN;
                    b.writeUnsignedInt(1)
                    b.writeUTFBytes(JSON.stringify(this.LoginHall));
                    this.send(b)
                break;
            case this.croomBtn:
                    var b: egret.ByteArray = new egret.ByteArray();
                    b.endian = egret.Endian.LITTLE_ENDIAN;
                    b.writeUnsignedInt(2)
                    // b.writeUTFBytes(JSON.stringify(this.LoginHall));
                    this.send(b)
                break;
            case this.joinBtn:
                    var b: egret.ByteArray = new egret.ByteArray();
                    b.endian = egret.Endian.LITTLE_ENDIAN;
                    b.writeUnsignedInt(3)
                    this.Croom.rid=parseInt(this.roomEdit.text)
                    b.writeUTFBytes(JSON.stringify(this.Croom));
                    this.send(b)

                break;
            case this.rinfoBtn:
                break;
        }

    }

        private socket:egret.WebSocket;

    private createSocket(){

        this.socket = new egret.WebSocket();
        this.socket.type = egret.WebSocket.TYPE_BINARY;
        this.socket.addEventListener(egret.Event.CONNECT, this.onConnect, this);         
        this.socket.addEventListener(egret.Event.CLOSE,this.onClose,this);
        this.socket.addEventListener(egret.IOErrorEvent.IO_ERROR,this.onError,this);
        this.socket.addEventListener(egret.ProgressEvent.SOCKET_DATA,this.onRecieve,this);
        this.socket.connectByUrl("ws://127.0.0.1:8080");
    }

    private onClose(e:egret.Event):void {
        console.log("close")
    }

    //连接错误
    private onError(e:egret.IOErrorEvent):void {

    }
        //接收数据
    private onRecieve(e: egret.ProgressEvent): void {
        // console.log("--------------------------------------")
        var b: egret.ByteArray = new egret.ByteArray();
        // b.endian = egret.Endian.LITTLE_ENDIAN;
        this.socket.readBytes(b);       
        console.log(b.length) 
        // var echo=this.socket.readUTF();
        var cmd=b.readUnsignedInt();
        var data=b.readUTFBytes(b.length-4);
        console.log(cmd,data)
        // console.log(echo);
        // this.socket.close();
    }

       //连接成功
    private onConnect(e:egret.Event):void {
        egret.log(this.name+ " connect success");
    }

    private send(data){
        this.socket.writeBytes(data);
        this.socket.flush();

    }
}