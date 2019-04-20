import { Component, AfterViewInit } from '@angular/core';
import { Router} from '@angular/router';

@Component({
  selector: 'app-admin-welcome',
  templateUrl: './admin-welcome.component.html',
  styleUrls: ['./admin-welcome.component.css']
})
export class AdminWelcomeComponent implements AfterViewInit {
	stateData
	constructor(private router: Router) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
	}

	ngAfterViewInit() {

	}

	companiesRedirect(event: any) {
		// Append companyName to stateData
		this.router.navigateByUrl('/admin/companies', { state: this.stateData })
	}

	statsRedirect(event: any) {
		// Append companyName to stateData
		this.router.navigateByUrl('/admin/stats', { state: this.stateData })
	}
}
