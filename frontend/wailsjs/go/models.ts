export namespace config {
	
	export class Item {
	    provider: string;
	    name: string;
	    fileName: string;
	    path: string;
	    scope: string;
	    format: string;
	    exists: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.name = source["name"];
	        this.fileName = source["fileName"];
	        this.path = source["path"];
	        this.scope = source["scope"];
	        this.format = source["format"];
	        this.exists = source["exists"];
	    }
	}

}

export namespace marketplaces {
	
	export class MCPPackage {
	    name: string;
	    description: string;
	    version: string;
	    author: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new MCPPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.version = source["version"];
	        this.author = source["author"];
	        this.url = source["url"];
	    }
	}

}

