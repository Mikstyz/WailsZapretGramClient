export namespace ethernet {
	
	export class TcpClient {
	    Conn: any;
	    PendingReqs: Record<string, >;
	    IP: string;
	    Port: string;
	    // Go type: Tools
	    Key?: any;
	    UserId: number;
	    Name: string;
	    Token: string;
	    Chats: Record<string, model.Chat>;
	
	    static createFrom(source: any = {}) {
	        return new TcpClient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Conn = source["Conn"];
	        this.PendingReqs = source["PendingReqs"];
	        this.IP = source["IP"];
	        this.Port = source["Port"];
	        this.Key = this.convertValues(source["Key"], null);
	        this.UserId = source["UserId"];
	        this.Name = source["Name"];
	        this.Token = source["Token"];
	        this.Chats = this.convertValues(source["Chats"], model.Chat, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TcpRequest {
	    Tcp?: TcpClient;
	
	    static createFrom(source: any = {}) {
	        return new TcpRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Tcp = this.convertValues(source["Tcp"], TcpClient);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace model {
	
	export class Chat {
	    Id: number;
	
	    static createFrom(source: any = {}) {
	        return new Chat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	    }
	}

}

