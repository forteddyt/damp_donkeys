import { Component, OnInit } from '@angular/core';
import { Router} from '@angular/router';

@Component({
	selector: 'app-add-company',
	templateUrl: './add-company.component.html',
	styleUrls: ['./add-company.component.css']
})
export class AddCompanyComponent implements OnInit {
	stateData
	constructor(private router: Router) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
	}
	ngOnInit() {
	}

	companiesRedirect(event: any){
		this.router.navigateByUrl('/admin/companies', { state: this.stateData })
	}
}
