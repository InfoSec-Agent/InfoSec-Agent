export namespace checks {
	
	export class Check {
	    issue_id: number;
	    result_id: number;
	    result?: string[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Check(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issue_id = source["issue_id"];
	        this.result_id = source["result_id"];
	        this.result = source["result"];
	        this.error = source["error"];
	    }
	}

}

export namespace usersettings {
	
	export class UserSettings {
	    Language: number;
	    ScanInterval: number;
	    Integration: boolean;
	    // Go type: time
	    NextScan: any;
	    Points: number;
	    PointsHistory: number[];
	    TimeStamps: time.Time[];
	    LighthouseState: number;
	    ProgressBarState: number;
	    IntegrationKey: string;
	
	    static createFrom(source: any = {}) {
	        return new UserSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Language = source["Language"];
	        this.ScanInterval = source["ScanInterval"];
	        this.Integration = source["Integration"];
	        this.NextScan = this.convertValues(source["NextScan"], null);
	        this.Points = source["Points"];
	        this.PointsHistory = source["PointsHistory"];
	        this.TimeStamps = this.convertValues(source["TimeStamps"], time.Time);
	        this.LighthouseState = source["LighthouseState"];
	        this.ProgressBarState = source["ProgressBarState"];
	        this.IntegrationKey = source["IntegrationKey"];
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

