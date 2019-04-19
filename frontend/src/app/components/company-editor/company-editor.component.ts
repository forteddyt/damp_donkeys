import { Component, OnInit } from '@angular/core';
import { Router} from '@angular/router';

@Component({
  selector: 'app-company-editor',
  templateUrl: './company-editor.component.html',
  styleUrls: ['./company-editor.component.css']
})
export class CompanyEditorComponent {
	stateData
	constructor(private router: Router) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
	}

	addRedirect(event: any) {
		// Append companyName to stateData
		this.router.navigateByUrl('/admin/companies/add', { state: this.stateData })
	}

	viewRedirect(event: any){
		// Add clicked on company to stateData
		let temp = "temp" // would be whatever user clicked on
		this.router.navigateByUrl('/admin/companies/view/' + temp, { state: this.stateData })
	}

	backRedirect(event: any){
		this.router.navigateByUrl('/admin/nav', { state: this.stateData })
	}
}
