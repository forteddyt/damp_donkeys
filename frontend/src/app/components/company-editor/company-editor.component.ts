import { Component, AfterViewInit, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { Router} from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { CompanyNameTileComponent } from '../company-name-tile/company-name-tile.component';

@Component({
  selector: 'app-company-editor',
  templateUrl: './company-editor.component.html',
  styleUrls: ['./company-editor.component.css']
})
export class CompanyEditorComponent implements AfterViewInit {
	@ViewChild('companyInsert', { read: ViewContainerRef }) companyInsert: ViewContainerRef;
	stateData
	companyList
	newCompanyName

	constructor(private http: HttpClient, private router: Router, private vcr: ViewContainerRef, private cfr: ComponentFactoryResolver) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
	}

	async ngAfterViewInit() {
		const resp = await this.getCompanies();

		//console.log(resp.status)
		this.companyList = resp.body["company_list"];
		this.loadComponents()
	}

	getCompanies() {
		return this.http.get("https://csrcint.cs.vt.edu/api/company_list", {observe: 'response'}).toPromise();
	}

	loadComponents() {
		//console.log(this.companyList)
		const cFactory = this.cfr.resolveComponentFactory(CompanyNameTileComponent);
		this.companyInsert.clear();
		for(var i in this.companyList){
			const companyComponent = <CompanyNameTileComponent>this.companyInsert.createComponent(cFactory).instance;

			companyComponent.companyName = this.companyList[i];
			companyComponent.stateData = this.stateData;
		}
	}

	onKey(event: any)
	{
		this.newCompanyName = event.target.value;
	}

	addCompany(event: any){
		console.log(this.newCompanyName)
		if (this.newCompanyName == null || this.newCompanyName == ""){
			alert("Invalid company name")
		} else {
			this.addCompanyToDB().then(
				(val) => { // success
					console.log(val)
					this.companyList.push(this.newCompanyName)
					// resort, ignoring case
					this.companyList.sort(function (a, b){
						return a.toLowerCase().localeCompare(b.toLowerCase());
					})
					this.stateData.jwt = val["jwt"] // update jwt
					alert("Company \"" + val["company_name"] + "\" has login code: " + val["user_code"])
					this.loadComponents()
				},
				(err) => { // failure
					if (err.status == 401){ // invalid jwt for request
						this.router.navigateByUrl('/admin')
					} else {
						alert("Whoops, something broke... status code: " + err.status)
					}
				});
		}
		// this.router.navigateByUrl('/admin/companies', { state: this.stateData })
	}

	addCompanyToDB() {
		return this.http.put("https://csrcint.cs.vt.edu/api/add_company?company_name=" + this.newCompanyName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
	}

	backRedirect(event: any){
		this.router.navigateByUrl('/admin/nav', { state: this.stateData })
	}
}
