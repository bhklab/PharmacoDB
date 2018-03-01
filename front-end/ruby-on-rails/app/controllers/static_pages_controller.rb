require 'csv'

class StaticPagesController < ApplicationController

    skip_before_action :verify_authenticity_token

    def home
        if not cookies.key?(:visited)
            cookies[:visited] = {
                :value => 0,
                :expires => 1.year.from_now
            }
        else
            # cookies[:visited] = {
            # 	:value => true,
            # 	:expires => 1.year.from_now
            # }
        end

        @cell_lines_count = Cell.count
        @tissues_count = Tissue.count
        @drugs_count = Drug.count
        @datasets_count = Dataset.count
        @experiments_count = Experiment.count
        @genes_count = Gene.count

        if params[:title].present? && params[:message].present?
            @title = params[:title]
            @message = params[:message]

            @name = params[:name].present? ? params[:name] : "n/a"
            @email = params[:email].present? ? parconams[:email] : "n/a"
            @subject = params[:subject].present? ? params[:subject] : "n/a"
            @id = params[:id].present? ? params[:id] : "n/a"

            Github.configure do |c|
                c.basic_auth = 'user:pw'
                c.user       = 'user'
                c.repo       = 'repo'
            end

            if @subject != "n/a"
                url = "/%s" % [@subject.gsub('-','_')]
                if @subject != 'experiments'
                    url = "%s/%s" % [url, @id]
                end
                @message = "%s\n\nURL: %s%s" % [@message, "https://pharmacodb.pmgenomics.ca", url]
                redirect_to url
            end

            github = Github.new
            github.issues.create title: @title, body: @message, labels: [@subject]    


        end
    end

    def cite_us

    end

    def download
        if params[:cell_annotation]
            send_file(Rails.root.join('app' , 'assets', 'csv', 'cell_annotation_table.csv'), :filename => 'cell_annotation_table-' + DateTime.now.to_s + '.csv')			
        end
        if params[:drug_annotation]
            send_file(Rails.root.join('app' , 'assets', 'csv', 'drug_annotation_table.csv'), :filename => 'drug_annotation_table-' + DateTime.now.to_s + '.csv')			
        end
        @version = ""
        if !params[:v].blank?
            @version = params[:v]
        end
    end

    def batch_query
        cell_names = params[:cell_names]
        drug_names = params[:drug_names]
        if cell_names.present? && drug_names.present?
            base_query = "SELECT cell_name, drug_name, dataset_name, HS, Einf, EC50, AAC, IC50, DSS1, DSS2, DSS3 FROM `experiments` " \
            "INNER JOIN `cells` ON `cells`.`cell_id` = `experiments`.`cell_id` " \
            "INNER JOIN `drugs` ON `drugs`.`drug_id` = `experiments`.`drug_id` " \
            "INNER JOIN `profiles` ON `profiles`.`experiment_id` = `experiments`.`experiment_id` " \
            "INNER JOIN `datasets` ON `datasets`.`dataset_id` = `experiments`.`dataset_id` "
            # check presence of wild star and build WHERE query
            cell_query = ""
            if !(cell_names.length == 1 && cell_names[0] == "*")
                cell_query = "cell_name IN (#{cell_names.map{ |e| "'" + e + "'" }.join(', ')})"
            end
            drug_query = ""
            if !(drug_names.length == 1 && drug_names[0] == "*")
                drug_query = "drug_name IN (#{drug_names.map{ |e| "'" + e + "'" }.join(', ')})"
            end
            where_query = "WHERE 1"
            if !cell_query.blank?
                where_query = where_query + " AND " + cell_query
            end
            if !drug_query.blank?
                where_query = where_query + " AND " + drug_query
            end
            sql = base_query + where_query
            results = ActiveRecord::Base.connection.exec_query(sql)
            if results.count > 0
                data_csv = ''
                data_csv << "cell name, drug name, dataset name, HS, Einf, EC50, AAC, IC50, DSS1, DSS2, DSS3\n"
                results.each do |r|
                    data_csv << r.values.join(', ')
                    data_csv << "\n"
                end
                send_data data_csv, filename: 'batch_query_results-' + DateTime.now.to_s + '.csv'
            else 
                respond_to do |format|
                    format.json { render :json => "1" }
                end
            end
        end
    end

    def news
        @version = ""
        if !params[:v].blank?
            @version = params[:v]
        end
    end

    def pharmacogx
        @dataset_id = []
        @dataset_name = []
        if params[:pgx].present?
            dataset = Dataset.find(params[:pgx].to_i)
            @dataset_id << dataset.dataset_id
            @dataset_name << dataset.dataset_name
        end
    end

end
