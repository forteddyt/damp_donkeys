import { Component, OnInit } from '@angular/core';
import { Router} from '@angular/router';

@Component({
	selector: 'app-view-company',
	templateUrl: './view-company.component.html',
	styleUrls: ['./view-company.component.css']
})

export class ViewCompanyComponent implements OnInit {
	stateData
	companyName
	constructor(private router: Router) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null
			|| this.stateData.companyName == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
		this.companyName = this.stateData.companyName
	}

	ngOnInit() {
	}

	companiesRedirect(event: any){
		this.router.navigateByUrl('/admin/companies', { state: this.stateData })
	}

	backRedirect(event: any){
		this.router.navigateByUrl('/admin/companies', { state: this.stateData })
	}

	generateCode(event: any){
		// Make api call to code generator api here
		alert("New code is: ")
	}
}
