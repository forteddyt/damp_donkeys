import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

export interface editInterface {
    selectCareerFairHelper(careerFairName: string, newCompanyList: string[], newJWT: string);
}

@Component({
  selector: 'app-career-fair-tile',
  templateUrl: './career-fair-tile.component.html',
  styleUrls: ['./career-fair-tile.component.css']
})
export class CareerFairTileComponent implements OnInit {
  careerFairName
  stateData

  public selfRef: CareerFairTileComponent

  //interface for Parent-Child interaction
  public compInteraction: editInterface;


  constructor(private http: HttpClient, private router: Router) { }

  ngOnInit() {
  }

  switchCareerFairs(){
  	this.companyListHelper().then(
  		(resp) => {
	  		this.compInteraction.selectCareerFairHelper(this.careerFairName,resp.body["company_list"], resp.body["jwt"])
  		},
  		(err) => {
  			if (err.status == 401){ // invalid jwt for request
				alert("Session expired")
				this.router.navigateByUrl('/admin')
  			} else {
				alert("Whoops, something broke... status code: " + err.status)
  			}
  		});
  }

  companyListHelper(){
  	return this.http.get("https://csrcint.cs.vt.edu/api/company_list?career_fair_name=" + this.careerFairName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
  }
}
