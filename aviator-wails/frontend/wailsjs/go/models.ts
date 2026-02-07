export namespace config {
	
	export class App {
	    id: string;
	    name: string;
	    path: string;
	    args: string;
	
	    static createFrom(source: any = {}) {
	        return new App(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.args = source["args"];
	    }
	}

}

