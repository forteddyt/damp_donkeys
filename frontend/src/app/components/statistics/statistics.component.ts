import { Component, OnInit } from '@angular/core';
import { Router} from '@angular/router';

@Component({
	selector: 'app-statistics',
	templateUrl: './statistics.component.html',
	styleUrls: ['./statistics.component.css']
})
export class StatisticsComponent implements OnInit {
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

	viewRedirect(event: any){
		this.router.navigateByUrl('/admin/nav', { state: this.stateData })
	}
}
