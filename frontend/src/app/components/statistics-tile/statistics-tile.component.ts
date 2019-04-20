import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { interval, Subscription } from 'rxjs';


@Component({
	selector: 'app-statistics-tile',
	templateUrl: './statistics-tile.component.html',
	styleUrls: ['./statistics-tile.component.css']
})
export class StatisticsTileComponent implements OnInit {
	sub: Subscription;

	stateData
	careerFairName
	companies
	numStudents
	numInterviews

	constructor(private http: HttpClient) {}

	ngOnInit() {
		const source = interval(60 * 1000);
		this.updateStats();
		this.sub = source.subscribe((val) => { this.updateStats() })
	}

	ngOnDestroy(){
		this.sub && this.sub.unsubscribe();
	}

	getStats() {
		return this.http.get("https://csrcint.cs.vt.edu/api/career_fair_stats?career_fair_name=" + this.careerFairName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
	}

	async updateStats() {
		console.log("updating...")
		let resp = await this.getStats()
		if (resp.status != 200) {
			this.companies = "error!"
			this.numStudents = "error!"
			this.numInterviews = "error!"
			return
		}

		this.companies = resp.body["company_list"]
		this.numStudents = resp.body["student_count"]
		this.numInterviews = resp.body["interview_count"]
	}
}
