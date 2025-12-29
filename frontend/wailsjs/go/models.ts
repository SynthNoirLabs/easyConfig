export namespace config {
	
	export class DocsPage {
	    provider: string;
	    title: string;
	    slug: string;
	    date: string;
	    hasMarkdown: boolean;
	    hasHtml: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DocsPage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.title = source["title"];
	        this.slug = source["slug"];
	        this.date = source["date"];
	        this.hasMarkdown = source["hasMarkdown"];
	        this.hasHtml = source["hasHtml"];
	    }
	}
	export class DocsProvider {
	    provider: string;
	    date: string;
	    pages: DocsPage[];
	
	    static createFrom(source: any = {}) {
	        return new DocsProvider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.date = source["date"];
	        this.pages = this.convertValues(source["pages"], DocsPage);
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
	export class ProfileSummary {
	    name: string;
	    itemCount: number;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ProfileSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.itemCount = source["itemCount"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class ProviderStatus {
	    providerName: string;
	    health: string;
	    statusMessage?: string;
	    discoveredFiles?: Item[];
	    lastChecked: string;
	
	    static createFrom(source: any = {}) {
	        return new ProviderStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.providerName = source["providerName"];
	        this.health = source["health"];
	        this.statusMessage = source["statusMessage"];
	        this.discoveredFiles = this.convertValues(source["discoveredFiles"], Item);
	        this.lastChecked = source["lastChecked"];
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

export namespace main {

	export class MarketplaceCacheStatus {
	    isCached: boolean;
	    isStale: boolean;

	    static createFrom(source: any = {}) {
	        return new MarketplaceCacheStatus(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isCached = source["isCached"];
	        this.isStale = source["isStale"];
	    }
	}

}

export namespace marketplaces {
	
	export class MCPPackage {
	    name: string;
	    description: string;
	    vendor?: string;
	    source: string;
	    url?: string;
	    version?: string;
	    author?: string;
	    stars?: number;
	    downloads?: number;
	    tags?: string[];
	    repoUrl?: string;
	    license?: string;
	    verified?: boolean;
	    checksum?: string;
	
	    static createFrom(source: any = {}) {
	        return new MCPPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.vendor = source["vendor"];
	        this.source = source["source"];
	        this.url = source["url"];
	        this.version = source["version"];
	        this.author = source["author"];
	        this.stars = source["stars"];
	        this.downloads = source["downloads"];
	        this.tags = source["tags"];
	        this.repoUrl = source["repoUrl"];
	        this.license = source["license"];
	        this.verified = source["verified"];
	        this.checksum = source["checksum"];
	    }
	}

}

export namespace workflows {
	
	export class Template {
	    id: string;
	    name: string;
	    description: string;
	    agent: string;
	    trigger: string;
	    tags: string[];
	    defaultFilename: string;
	    content: string;
	    requiredSecrets: string[];
	    setupInstructions: string;
	
	    static createFrom(source: any = {}) {
	        return new Template(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.agent = source["agent"];
	        this.trigger = source["trigger"];
	        this.tags = source["tags"];
	        this.defaultFilename = source["defaultFilename"];
	        this.content = source["content"];
	        this.requiredSecrets = source["requiredSecrets"];
	        this.setupInstructions = source["setupInstructions"];
	    }
	}
	export class WorkflowResponse {
	    content: string;
	    requiredSecrets: string[];
	    setupInstructions: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.requiredSecrets = source["requiredSecrets"];
	        this.setupInstructions = source["setupInstructions"];
	    }
	}

}

