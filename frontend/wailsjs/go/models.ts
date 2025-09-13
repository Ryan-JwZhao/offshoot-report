export namespace main {
	
	export class ReportRequest {
	    projectTitle: string;
	    backups: string;
	    filePaths: string[];
	
	    static createFrom(source: any = {}) {
	        return new ReportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectTitle = source["projectTitle"];
	        this.backups = source["backups"];
	        this.filePaths = source["filePaths"];
	    }
	}

}

