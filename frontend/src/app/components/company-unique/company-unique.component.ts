import { Component, AfterViewInit, ComponentFactoryResolver, ViewChild, ViewContainerRef, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import * as jwt_decode from 'jwt-decode';
import { InterviewComponent } from "../interview/interview.component"

@Component({
  selector: 'app-company-unique',
  templateUrl: './company-unique.component.html',
  styleUrls: ['./company-unique.component.css']
})
export class CompanyUniqueComponent implements OnInit {
  @ViewChild('interview', { read: ViewContainerRef }) companyInsert: ViewContainerRef;
  code = "";
  validString: String = "abcdefghij";
  company_name = "Company_Name";
  json
  //json_middle
  api_call_return = {"company_name": "test_company", "interviewees": [{"name": "Joe Smith", "check-in": "11:30"}, {"name": "Jill Smith", "check-in": "12:00"}]};
  //name_array = ["Joe Smith", "Jill Smith"];
  //time_arry = ["11:30", "12:00"];
  companyList

  constructor(private http: HttpClient, private route: ActivatedRoute, private router: Router, private cfr: ComponentFactoryResolver) {
     //code = this.route.snapshot.paramMap.get('code');
  }

  ngOnInit() {

    //console.log(this.api_call_return['interviewees'][0]['name']);
    //console.log(this.json_middle);
    //console.log(jwt_decode(this.json));

    //this.company_name = this.json.jwt;

    //this.code = this.route.snapshot.paramMap.get('code');
    //this.code = this.route.snapshot.params.code;
    //if (angular.equals(this.code, this.validString)) {}
    //else {
    //  this.router.navigate(['employers'], {});
    //}
  }

  getUser()
  {
    return this.http.get("https://csrcint.cs.vt.edu/api/login?code=" + this.code, {observe: 'response'}).toPromise();
  }

  //login()
  //{
    //return this.http.get("https://csrcint.cs.vt.edu/api/login?password_hash=test", {observe: 'response'}).toPromise();
      //console.log(res);
      //var res = await this.getUser();
      //this.json = res.body["jwt"];
      //this.json = res["jwt"];
      //this.json_middle = this.json.split('.')[1];
      //var decoded = jwt_decode(this.json);
      //console.log(decoded);
      //this.company_name = decoded['user'];
  //}

  async ngAfterViewInit() {

    this.code = this.route.snapshot.params.code;
    //console.log(this.code);
    //const resp = /*await*/ this.getNames();

    //console.log(resp.status)
    //this.companyList = resp['interviewees'];//.body
    //console.log(this.companyList);
    //this.loadComponents();
    //while (this.json == undefined) {}

    //var res = await this.getUser();

    //console.log(this.json);
    const r = await this.getUser();

    this.json = r.body['jwt'];
    var decoded = jwt_decode(this.json);
    this.company_name = decoded['user'];

    const resp = await this.getNames();

    //console.log(resp.status);
    this.companyList = resp.body['students'];
    this.loadComponents();
  }

  getNames()
  {
    return this.http.get("https://csrcint.cs.vt.edu/api/company_check_ins?company_name=" + this.company_name + "&jwt=" + this.json, {observe: 'response'}).toPromise();
    // update jwt
  }

  loadComponents() {
    //console.log(this.companyList);
    const cFactory = this.cfr.resolveComponentFactory(InterviewComponent);
    //this.companyInsert.clear();
    for(var i in this.companyList){
      //console.log("for loop executed");
      const companyComponent = <InterviewComponent>this.companyInsert.createComponent(cFactory).instance;

      companyComponent.name = this.companyList[i]['name'];
      companyComponent.time = this.companyList[i]['check-in'];
    }
  }
}
