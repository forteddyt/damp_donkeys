import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
	selector: 'app-statistics-tile',
	templateUrl: './statistics-tile.component.html',
	styleUrls: ['./statistics-tile.component.css']
})
export class StatisticsTileComponent implements OnInit {
	stateData
	careerFairName
	companies
	numStudents
	numInterviews

	constructor(private http: HttpClient) {}

	async ngOnInit() {
		let resp = await this.getStats()
		if (resp.status != 200) {
			this.companies = "!"
			this.numStudents = "!"
			this.numInterviews = "!"
			return
		}

		this.companies = resp.body["company_list"]
		this.numStudents = resp.body["student_count"]
		this.numInterviews = resp.body["interview_count"]
	}

	getStats() {
		return this.http.get("https://csrcint.cs.vt.edu/api/career_fair_stats?career_fair_name=" + this.careerFairName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
	}
}
