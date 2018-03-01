require 'json'

class ProfilesController < ApplicationController

  skip_before_action :verify_authenticity_token

   def is_equal_incl(arg, name)
    return (arg.casecmp(name) == 0) || (name.to_s.downcase().include? arg.downcase())
   end

   def docs
     @toc = [ ['Cell Lines', 0], ['Cell Lines vs Drugs', 1], ['Datasets', 0],
              ['Datasets Intersection', 1], ['Drugs', 0], ['Search', 0],
              ['Tissues', 0]
            ]
     @exit_code = 1
     return
   end

   def cell_lines
      unless params[:id].present?
         tissues = Tissue.all
         @count_ids = []
         @count_names = []
         tissues.each do |t|
            sql = "SELECT DISTINCT cell_id FROM cell_tissues WHERE tissue_id = #{t.tissue_id}"
            @count_ids << ActiveRecord::Base.connection.exec_query(sql).count
            @count_names << t.tissue_name
         end
         top_cutoff = @count_ids.sort.reverse[12]
         collapse_x = @count_ids.each_index.select{|x| @count_ids[x] <= top_cutoff}
         total_o = collapse_x.map{|x| @count_ids[x]}.sum
         @other_names = collapse_x.map{|x| @count_names[x]}
         @other_nums = collapse_x.map{|x| @count_ids[x]}
         not_collapse_x = @count_ids.each_index.select{|x| @count_ids[x] > top_cutoff}
         @count_ids = @count_ids.find_all{|x| x > top_cutoff}
         @count_names = not_collapse_x.map{|x| @count_names[x].html_safe}
         @count_ids << total_o
         @count_names << "Other"
         sql = "SELECT cell_id, cell_name FROM cells"
         cell_lines = ActiveRecord::Base.connection.exec_query(sql)
         filtered_cell_lines = []
         cell_lines.each do |cl|
           if params[:c].to_s.strip.empty?
             filtered_cell_lines << cl
           else
             filtered_cell_lines << cl if is_equal_incl(params[:c], cl['cell_name'])
           end
         end
         @cell_lines = filtered_cell_lines.to_a.paginate(:page => params[:page], :per_page => 20)
         @cell_lines_count = cell_lines.count
         @exit_code = 0
         return
      end

      sql = "SELECT cell_id, cell_name FROM cells WHERE cell_id = #{params[:id]}"
      cell_line = ActiveRecord::Base.connection.exec_query(sql)
      @cell_line_id = cell_line[0]['cell_id']
      @cell_line_name = cell_line[0]['cell_name']

      @synonyms = []
      @diseases = ""
      @ncit_path = "https://ncit.nci.nih.gov/ncitbrowser/"
      @cellosaurus_path = "http://web.expasy.org/cellosaurus/"

      sql = "SELECT accession, sy, di FROM cellosaurus WHERE identifier = '#{@cell_line_name}' OR sy LIKE '%#{@cell_line_name}%'"
      cellosaurus_data = ActiveRecord::Base.connection.exec_query(sql)
      if cellosaurus_data.present?
         accession = cellosaurus_data[0]['accession']
         @cellosaurus_path = @cellosaurus_path + accession
         if cellosaurus_data[0]['synonyms'].present?
            synonyms = cellosaurus_data[0]['synonyms'].split("; ").join(", ")
            html_link = "<a target=\"_blank\" href=\"http://web.expasy.org/cellosaurus/#{accession}\" target=\"_blank\">Cellosaurus</a>"
            @synonyms << [html_link, synonyms]
         end
         if cellosaurus_data[0]['di'].present?
            id = cellosaurus_data[0]['di'].split("; ")[1]
            @diseases = cellosaurus_data[0]['di'].split("; ").join(", ")
            @ncit_path = "https://ncit.nci.nih.gov/ncitbrowser/ConceptReport.jsp?dictionary=NCI%20Thesaurus&code=#{id}"
         else
            @diseases = "N/A"
         end
      end

      sql = "SELECT cell_name, dataset_id, dataset_name FROM source_cell_names, datasets WHERE cell_id = #{@cell_line_id} AND dataset_id = source_id"
      syn = ActiveRecord::Base.connection.exec_query(sql)
      if syn.present?
         allsyn = []
         syn.each do |s|
            dataset_id = s['dataset_id']
            dataset_name = s['dataset_name']
            html_link = "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
            allsyn << [html_link, s['cell_name']]
         end
         names_so_far = []
         sources_so_far = []
         allsyn.each do |a|
            if names_so_far.include? a[1]
               if sources_so_far.include? a[0]
                  next
               end
               @synonyms.each do |s|
                  if s[1] == a[1]
                     if s[0].include? a[0]
                        break
                     else
                        s[0] = s[0] + ", " + a[0]
                        sources_so_far << a[0]
                        break
                     end
                  end
               end
            elsif sources_so_far.include? a[0]
               @synonyms.each do |s|
                  if s[0] == a[0]
                     s[1] = s[1] + ", " + a[1]
                     names_so_far << a[1]
                     break
                  end
               end
               unless names_so_far.include? a[1]
                  @synonyms << a
                  sources_so_far << a[0]
                  names_so_far << a[1]
               end
            else
               @synonyms << a
               sources_so_far << a[0]
               names_so_far << a[1]
            end
         end
      end

      unless @synonyms.present?
         @synonyms = "N/A"
      end

      mol_data = MolCell.where(:cell_id => @cell_line_id).includes(:dataset).pluck(:dataset_name, :mDataType, :num_prof)

      @mol_data_type = mol_data.pluck(1).uniq

      mol_data_datasets = mol_data.pluck(0).uniq

      @mol_data = []

      for dataset in mol_data_datasets
         entry = []
         entry << dataset
         for type in @mol_data_type
            nn = mol_data.select{|x| x[0] == dataset}.select{|x| x[1] == type}.pluck(2).sum
            unless nn.nil?
               entry << nn
            else
               entry << 0
            end
         end
         @mol_data << entry
      end

      sql = "SELECT t.tissue_id, t.tissue_name FROM tissues t, cell_tissues ct WHERE ct.cell_id = #{@cell_line_id} AND t.tissue_id = ct.tissue_id"
      tissue = ActiveRecord::Base.connection.exec_query(sql)
      @tissue_id = tissue[0]['tissue_id']
      @tissue_name = tissue[0]['tissue_name']

      sql = "SELECT d.drug_id, d.drug_name, s.dataset_id, s.dataset_name, e.experiment_id, p.IC50, p.AAC FROM drugs d, datasets s, experiments e, profiles p WHERE e.cell_id = #{@cell_line_id} AND d.drug_id = e.drug_id AND s.dataset_id = e.dataset_id AND p.experiment_id = e.experiment_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)

      drugs = []
      datasets = []
      dataset_count = []
      drugs_arr = []
      @synonyms_waterfall = []

      @drug_AAC = []
      @drug_IC50 = []

      waterfall_data = []

      experiments.each do |e|
         drug_id = e['drug_id']
         drug_name = e['drug_name']
         dataset_id = e['dataset_id']
         dataset_name = e['dataset_name']
         if dataset_count.flatten.include? dataset_name
            dataset_count.each do |s|
               if s[0] == dataset_name
                  s[1] += 1
                  break
               end
            end
         else
            dataset_count << [dataset_name, 1]
         end
         if drugs.include? drug_name
            drugs_arr.each do |d|
               if d[0] == drug_name
                  unless d[1].include? dataset_name
                     d[1] = d[1] + ", " + "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
                  end
                  d[2] += 1
                  datasets << dataset_name
                  break
               end
            end
            ind = drugs.index(drug_name)
            tempAAC = e["AAC"]
            @drug_AAC[ind] << tempAAC
            if e["IC50"].nil?
              tempIC50 = Float::NAN
            else
              tempIC50 = e["IC50"]
            end
            @drug_IC50[ind] << tempIC50
            if params[:download] == "waterfall"
               waterfall_data << [drug_name, dataset_name, tempAAC, tempIC50]
            end
         else
            drugs << drug_name
            datasets << dataset_name
            tempAAC = e["AAC"]
            @drug_AAC << [tempAAC]

            if e["IC50"].nil?
              tempIC50 = Float::NAN
            else
              tempIC50 = e["IC50"]
            end
            @drug_IC50 << [tempIC50]
            if params[:d].to_s.strip.empty?
              drugs_arr << [drug_name, "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>", 1, drug_id]
            else
              drugs_arr << [drug_name, "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>", 1, drug_id] if is_equal_incl(params[:d], drug_name)
            end
            if params[:download] == "waterfall"
               waterfall_data << [dataset_name, drug_name, tempAAC, tempIC50]
            end
         end
         temp = [e["drug_name"]] + Drug.find(e['drug_id']).source_drug_names.pluck(:drug_name)
         @synonyms_waterfall << temp.uniq
   end

   # sql = "SELECT DISTINCT cell_name FROM source_cell_names WHERE cell_id IN (#{cell_ids.join(', ')})"
   # @synonyms_waterfall = []
   # results = ActiveRecord::Base.connection.exec_query(sql)
   # results.to_hash.each{ |k, v| @synonyms_waterfall.push(k.values[0]) }

      @synonyms_waterfall.uniq!

      @drug_names_waterfall = drugs

      #attempt at compiling synonyms for each drug of waterfall
     

      @numdatasets = datasets.uniq.count
      @numdrugs = drugs.count


      if params[:sort] == "compounds"
        drugs_arr = sort_column(drugs_arr, params[:direction], 0)
      elsif params[:sort] == "datasets"
        drugs_arr = sort_column(drugs_arr, params[:direction], 1)
      elsif params[:sort] == "experiments"
        drugs_arr = sort_column(drugs_arr, params[:direction], 2)
      else
        # sort the experiments row by default for both tables
        drugs_arr = sort_column(drugs_arr, "desc", 2)
      end

      @drugs = drugs_arr.paginate(:page => params[:page], :per_page => 10)

      if params[:download] == "drug_table"
         drugs_arr.each do |d|
            d[1] = d[1].gsub(/<a[^>]*>/, '')
            d[1] = d[1].gsub(/<\/a.?>/, '')
            d[1] = d[1].gsub(/,/, ';')
            d.delete_at(3)
         end
      end


      all_datasets = Dataset.pluck(:dataset_name)

      drugs_count = []

      sql = "SELECT d.drug_id, d.drug_name, s.dataset_id, s.dataset_name FROM drugs d, datasets s, experiments e WHERE e.cell_id = #{@cell_line_id} AND d.drug_id = e.drug_id AND s.dataset_id = e.dataset_id"
      exps = ActiveRecord::Base.connection.exec_query(sql)

      sql = "SELECT dataset_id, dataset_name FROM datasets"
      datasets = ActiveRecord::Base.connection.exec_query(sql)
      @ccounts = []
      datasets.each do |s|
         count = exps.select {|e| e['dataset_id'] == s['dataset_id']}.map{|x| x['drug_id']}.uniq.count
         drugs_count << [s['dataset_name'],count]
      end

      @snames = drugs_count.map { |x| x[0] }
      @scounts = drugs_count.map { |x| x[1] }

      # Set right flag for table caption

      if params[:d].to_s.strip.empty?
         @search_drugs = false
         @placeholder_drug = 'Search compound names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drug = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search compound names ...'
         else
            @caption_drug = ''
            @placeholder_drug = params[:d].to_s.strip
         end
      end

      if params[:download] == "drug_table"
         data_csv = ''
         data_csv << "Drug Name, Datasets, # Experiments\n"
         drugs_arr.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @cell_line_name + '_drug_table.csv'
      end

      if params[:download] == "waterfall"
         data_csv = ''
         data_csv << " Drug, Dataset, AAC, IC50\n"
         waterfall_data.sort_by!{|g| g[0]}
         waterfall_data.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @cell_line_name + '_waterfall_table.csv'
      end

      @exit_code = 1
      return
   end #/cell_lines

   def tissues
      unless params[:id].present?
                  tissues = Tissue.all
         @count_ids = []
         @count_names = []
         tissues.each do |t|
            sql = "SELECT DISTINCT cell_id FROM cell_tissues WHERE tissue_id = #{t.tissue_id}"
            @count_ids << ActiveRecord::Base.connection.exec_query(sql).count
            @count_names << t.tissue_name
         end
         top_cutoff = @count_ids.sort.reverse[12]
         collapse_x = @count_ids.each_index.select{|x| @count_ids[x] <= top_cutoff}
         total_o = collapse_x.map{|x| @count_ids[x]}.sum
         @other_names = collapse_x.map{|x| @count_names[x]}
         @other_nums = collapse_x.map{|x| @count_ids[x]}
         not_collapse_x = @count_ids.each_index.select{|x| @count_ids[x] > top_cutoff}
         @count_ids = @count_ids.find_all{|x| x > top_cutoff}
         @count_names = not_collapse_x.map{|x| @count_names[x]}
         @count_ids << total_o
         @count_names << "Other"

      sql = "SELECT tissue_id, tissue_name FROM tissues"
      tissues = ActiveRecord::Base.connection.exec_query(sql)
         filtered_tissues = []
         tissues.each do |ts|
           if params[:t].to_s.strip.empty?
             filtered_tissues << ts
           else
             filtered_tissues << ts if is_equal_incl(params[:t], ts['tissue_name'])
           end
         end
         @tissues = filtered_tissues.to_a.paginate(:page => params[:page], :per_page => 20)
         @tissues_count = tissues.count
         @exit_code = 0
         return
      end

      # sql = "SELECT tissue_name FROM tissues WHERE tissue_id = #{params[:id]}"
      # tissue = ActiveRecord::Base.connection.exec_query(sql)
      @tissue_id = params[:id]
      @tissue_name = Tissue.find(params[:id].to_i).tissue_name

      @synonyms = []

      sql = "SELECT tissue_name, dataset_id, dataset_name FROM source_tissue_names, datasets WHERE tissue_id = #{@tissue_id} AND dataset_id = source_id"
      syn = ActiveRecord::Base.connection.exec_query(sql)
      if syn.present?
         allsyn = []
         syn.each do |s|
            dataset_id = s['dataset_id']
            dataset_name = s['dataset_name']
            html_link = "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
            allsyn << [html_link, s['tissue_name'].tr("_", " ") ]
         end
         names_so_far = []
         sources_so_far = []
         allsyn.each do |a|
            if names_so_far.include? a[1]
               if sources_so_far.include? a[0]
                  next
               end
               @synonyms.each do |s|
                  if s[1] == a[1]
                     if s[0].include? a[0]
                        break
                     else
                        s[0] = s[0] + ", " + a[0]
                        sources_so_far << a[0]
                        break
                     end
                  end
               end
            elsif sources_so_far.include? a[0]
               @synonyms.each do |s|
                  if s[0] == a[0]
                     s[1] = s[1] + ", " + a[1]
                     names_so_far << a[1]
                     break
                  end
               end
               unless names_so_far.include? a[1]
                  @synonyms << a
                  sources_so_far << a[0]
                  names_so_far << a[1]
               end
            else
               @synonyms << a
               sources_so_far << a[0]
               names_so_far << a[1]
            end
         end
      end

      unless @synonyms.present?
         @synonyms = "N/A"
      end

      sql = "SELECT c.cell_id, c.cell_name FROM cells c, cell_tissues ct WHERE ct.tissue_id = #{@tissue_id} AND c.cell_id = ct.cell_id ORDER BY c.cell_name"
      cell_lines = ActiveRecord::Base.connection.exec_query(sql)
      filter_cell_lines = []
      cell_lines.to_a.each do |c|
        if params[:c].to_s.strip.empty?
          filter_cell_lines << c
        else
          filter_cell_lines << c if is_equal_incl(params[:c], c['cell_name'])
        end
      end

      @cell_lines = filter_cell_lines.paginate(:page => params[:cpage], :per_page => 20)
      @cell_lines_count = cell_lines.count

      array = []
      drugs = []
      datasets = []
      drugs_count = []
      cells_count = []

      sql = "SELECT d.drug_id, d.drug_name, s.dataset_id, s.dataset_name FROM drugs d, datasets s, experiments e, cell_tissues ct WHERE ct.tissue_id = #{@tissue_id} AND e.cell_id = ct.cell_id AND d.drug_id = e.drug_id AND s.dataset_id = e.dataset_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)
      ## unfortunately slow code with active record
      # experiments = Experiment.includes(:drug, :dataset).where(:tissue_id => @tissue_id)

      if experiments.present?
         experiments.each do |e|
            drug_id = e['drug_id']
            drug_name = e['drug_name']
            dataset_id = e['dataset_id']
            dataset_name = e['dataset_name']
            if drugs.include? drug_name
               array.each do |a|
                  if a[0] == drug_name
                     unless a[1].include? dataset_name
                        a[1] = a[1] + ", " + "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
                        datasets << dataset_name
                     end
                     a[2] += 1
                     break
                  end
               end
            else
              if params[:d].to_s.strip.empty?
                array << [drug_name, "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>", 1, drug_id]
              else
                array << [drug_name, "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>", 1, drug_id] if is_equal_incl(params[:d], drug_name)
              end
              drugs << drug_name
              datasets << dataset_name
            end
         end
      end

      if params[:sort] == "compounds"
        array = sort_column(array, params[:direction], 0)
      elsif params[:sort] == "datasets"
        array = sort_column(array, params[:direction], 1)
      elsif params[:sort] == "experiments"
        array = sort_column(array, params[:direction], 2)
      else
        # sort the experiments row by default for both tables
        array = sort_column(array, "desc", 2)
      end

      @numdrugs = drugs.uniq.count
      @numdatasets = datasets.uniq.count

      drugs_arr = array
      @drugs = array.paginate(:page => params[:dpage], :per_page => 10)

      if params[:download] == "drug_table"
         drugs_arr.each do |d|
            d[1] = d[1].gsub(/<a[^>]*>/, '')
            d[1] = d[1].gsub(/<\/a.?>/, '')
            d[1] = d[1].gsub(/,/, ';')
            d.delete_at(3)
         end
      end


      exps = experiments.to_a

      sql = "SELECT dataset_id, dataset_name FROM datasets"
      datasets = ActiveRecord::Base.connection.exec_query(sql)
      @ccounts = []
      datasets.each do |s|
         count = exps.select {|e| e['dataset_id'] == s['dataset_id']}.map{|x| x['drug_id']}.uniq.count
         drugs_count << [s['dataset_name'],count]
      end

      @snames = drugs_count.map { |x| x[0] }
      @scounts = drugs_count.map { |x| x[1] }


      sql = "SELECT ds.cell_id, ds.dataset_id FROM dataset_cells ds, cell_tissues ct WHERE ct.tissue_id = #{@tissue_id} AND ds.cell_id = ct.cell_id"
      dataset_cells = ActiveRecord::Base.connection.exec_query(sql)


      dss = dataset_cells.to_a
      @ccounts = []
      datasets.each do |s|
         count = dss.select {|e| e['dataset_id'] == s['dataset_id']}.map{|x| x['cell_id']}.uniq.count
         cells_count << [s['dataset_name'],count]
      end

      @ccounts = cells_count.map {|x| x[1] }

      # Set right flag for table caption

      if params[:d].to_s.strip.empty?
         @search_drugs = false
         @placeholder_drug = 'Search compound names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drug = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search compound names ...'
         else
            @caption_drug = ''
            @placeholder_drug = params[:d].to_s.strip
         end
      end



      if params[:c].to_s.strip.empty?
         @search_cell_lines = false
         @placeholder_cell_lines = 'Search cell line names ...'
      else
         @search_cell_lines = true
         if @cell_lines.length.equal? 0
            @caption_cell_lines = 'Your search for ' + params[:c].to_s.strip + " yielded no results."
            @placeholder_cell_lines = 'Search cell line names ...'
         else
            @caption_cell_lines = ''
            @placeholder_cell_lines = params[:c].to_s.strip
         end
      end



      if params[:download] == "cell_table"
         data_csv = ''
         data_csv << "Cell Name\n"
         filter_cell_lines.each do |d|
            data_csv << d["cell_name"]
            data_csv << "\n"
         end
         send_data data_csv, filename: @tissue_name + '_cell_table.csv'
      end


      if params[:download] == "drug_table"
         data_csv = ''
         data_csv << "Drug Name, Datasets, # Experiments\n"
         drugs_arr.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @tissue_name + '_drug_table.csv'
      end



      @exit_code = 1
      return
   end #/tissues

   def drugs
      @drugs_colspan = 2    
      @placeholder_drug = 'Search compound names ...'    
      unless params[:id].present?
         datasets = Dataset.all
         @count_ids = []
         @count_names = []
         datasets.each do |d|
            sql = "SELECT DISTINCT drug_id FROM experiments WHERE dataset_id = #{d.dataset_id}"
            @count_ids << ActiveRecord::Base.connection.exec_query(sql).count
            @count_names << d.dataset_name
         end

         sql = "SELECT drug_id, drug_name FROM drugs"
         drugs = ActiveRecord::Base.connection.exec_query(sql)
         filtered_drugs = []
         drugs.each do |drug|
           if params[:d].to_s.strip.empty?
             filtered_drugs << drug
           else
            filtered_drugs << drug if is_equal_incl(params[:d], drug['drug_name'])
           end
         end
         # XXX :: The database was changes as a workaround to get nice drugs first. Please fix the actual data!
         @drugs = filtered_drugs.sort_by{ |d| d['drug_id']}.to_a.paginate(:page => params[:page], :per_page => 20)
         @drugs_count = drugs.count
         if @drugs_count <= 10
            @drugs_colspan = 1
         end
         @exit_code = 0
         return
      end

      # sql = "SELECT drug_name FROM drugs WHERE drug_id = #{params[:id]}"
      drug = Drug.find(params[:id].to_i)
      @drug_id = params[:id]
      @drug_name = drug.drug_name

      @targets = drug.targets.uniq.to_a

      @synonyms = []
      @drug_ids = DrugAnnot.where(:drug_id => @drug_id).take
      @drug_ids = [["Smiles:",@drug_ids.smiles], ["Inchikey:",@drug_ids.inchikey], ["Pubchem:",@drug_ids.pubchem]]

      sql = "SELECT drug_name, dataset_id, dataset_name FROM source_drug_names, datasets WHERE drug_id = #{@drug_id} AND dataset_id = source_id"
      syn = ActiveRecord::Base.connection.exec_query(sql)
      if syn.present?
         allsyn = []
         syn.each do |s|
            dataset_id = s['dataset_id']
            dataset_name = s['dataset_name']
            html_link = "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
            allsyn << [html_link, s['drug_name']]
         end
         names_so_far = []
         sources_so_far = []
         allsyn.each do |a|
            if names_so_far.include? a[1]
               if sources_so_far.include? a[0]
                  next
               end
               @synonyms.each do |s|
                  if s[1] == a[1]
                     if s[0].include? a[0]
                        break
                     else
                        s[0] = s[0] + ", " + a[0]
                        sources_so_far << a[0]
                        break
                     end
                  end
               end
            elsif sources_so_far.include? a[0]
               @synonyms.each do |s|
                  if s[0] == a[0]
                     s[1] = s[1] + ", " + a[1]
                     names_so_far << a[1]
                     break
                  end
               end
               unless names_so_far.include? a[1]
                  @synonyms << a
                  sources_so_far << a[0]
                  names_so_far << a[1]
               end
            else
               @synonyms << a
               sources_so_far << a[0]
               names_so_far << a[1]
            end
         end
      end

      unless @synonyms.present?
         @synonyms = "N/A"
      end

      cell_line_array = []
      tissue_array = []
      cell_lines = []
      cell_ids = []
      tissues = []
      datasets = []

      @cell_AAC = []
      @cell_IC50 = []

      waterfall_data = []
      @synonyms_waterfall = []
      
      ### FIXME:: This following loop takes way too long because of the source cell name call at the end. How can we optimize the query?
      sql = "SELECT c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, s.dataset_id, s.dataset_name, e.experiment_id, p.AAC, p.IC50 FROM cells c, tissues t, datasets s, cell_tissues ct, experiments e, profiles p WHERE e.drug_id = #{@drug_id} AND s.dataset_id = e.dataset_id AND ct.cell_id = e.cell_id AND t.tissue_id = ct.tissue_id AND c.cell_id = e.cell_id AND p.experiment_id = e.experiment_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)

      if experiments.present?
         experiments.each do |e|
            cell_id = e['cell_id']
            cell_ids << cell_id
            cell_line = e['cell_name']
            tissue_id = e['tissue_id']
            # tissue = e['tissue_name']
            tissue = strip_underscore(e['tissue_name'])
            dataset_id = e['dataset_id']
            dataset = e['dataset_name']
            if cell_lines.include? cell_line
              cell_line_array.each do |a|
                  if a[0] == cell_line
                     unless a[1].include? dataset
                        a[1] = a[1] + ", " + "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>"
                        datasets << dataset
                     end
                     a[2] += 1
                     break
                  end
               end
               ind = cell_lines.index(cell_line)
               tempAAC = e["AAC"]
               @cell_AAC[ind] << tempAAC
               if e["IC50"].nil?
                 tempIC50 = Float::NAN
               else
                 tempIC50 = e["IC50"]
               end
               @cell_IC50[ind] << tempIC50
               if params[:download] == "waterfall"
                  waterfall_data << [cell_line, dataset, tempAAC, tempIC50]
               end
            else
              tempAAC = e["AAC"]
              @cell_AAC << [tempAAC]

              if e["IC50"].nil?
                tempIC50 = Float::NAN
              else
                tempIC50 = e["IC50"]
              end
              @cell_IC50 << [tempIC50]
              if params[:c].to_s.strip.empty?
                cell_line_array << [cell_line, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, cell_id, tissue]
              else
                cell_line_array << [cell_line, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, cell_id, tissue] if is_equal_incl(params[:c], cell_line)
              end
              cell_lines << cell_line
              datasets << dataset
              if params[:download] == "waterfall"
                  waterfall_data << [cell_line, dataset, tempAAC, tempIC50]
               end
            end
            if tissues.include? tissue
               tissue_array.each do |t|
                  if t[0] == tissue
                     unless t[1].include? dataset
                        t[1] = t[1] + ", " + "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>"
                        datasets << dataset
                     end
                     t[2] += 1
                     break
                  end
               end
            else
              if params[:t].to_s.strip.empty?
                tissue_array << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id]
              else
                tissue_array << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id] if is_equal_incl(params[:t], tissue)
              end
              tissues << tissue
              datasets << dataset
            end
            temp = [e["cell_name"]] + Cell.find(e['cell_id']).source_cell_names.pluck(:cell_name)
            @synonyms_waterfall << temp.uniq
         end
      end

      # sql = "SELECT DISTINCT cell_name FROM source_cell_names WHERE cell_id IN (#{cell_ids.join(', ')})"
      # @synonyms_waterfall = []
      # results = ActiveRecord::Base.connection.exec_query(sql)
      # results.to_hash.each{ |k, v| @synonyms_waterfall.push(k.values[0]) }

      @synonyms_waterfall.uniq!

      @cell_lines_waterfall = cell_lines

      @num_cell_lines = cell_lines.uniq.count
      @num_tissues = tissues.uniq.count
      @numdatasets = datasets.uniq.count
     
      # rearrange cell line and tissue tables
      if params[:sort] == "cell_line"
        cell_line_array = sort_column(cell_line_array, params[:direction], 0)
      elsif params[:sort] == "tissue"
        tissue_array = sort_column(tissue_array, params[:direction], 0)
      elsif params[:sort] == "tissue_type"
        cell_line_array = sort_column(cell_line_array, params[:direction], 4)
      elsif params[:sort] == "c_datasets"
        cell_line_array = sort_column(cell_line_array, params[:direction], 1)
      elsif params[:sort] == "t_datasets"
        tissue_array = sort_column(tissue_array, params[:direction], 1)
      elsif params[:sort] == "c_experiments"
        cell_line_array = sort_column(cell_line_array, params[:direction], 2)
      elsif params[:sort] == "t_experiments"
        tissue_array = sort_column(tissue_array, params[:direction], 2)
      else 
        # sort the experiments row by default for both tables
        cell_line_array = sort_column(cell_line_array, "desc", 2)
        tissue_array = sort_column(tissue_array, "desc", 2)
      end

      @cell_lines = cell_line_array.paginate(:page => params[:cpage], :per_page => 10)
      @tissues = tissue_array.paginate(:page => params[:tpage], :per_page => 10)

      if params[:download] == "cell_table"
         cell_line_array.each do |d|
            d[1] = d[1].gsub(/<a[^>]*>/, '')
            d[1] = d[1].gsub(/<\/a.?>/, '')
            d[1] = d[1].gsub(/,/, ';')
            d.delete_at(3)
         end
      end

      if params[:download] == "tissue_table"
         tissue_array.each do |d|
            d[1] = d[1].gsub(/<a[^>]*>/, '')
            d[1] = d[1].gsub(/<\/a.?>/, '')
            d[1] = d[1].gsub(/,/, ';')
            d.delete_at(3)
         end
      end

      @ccounts = []
      @tcounts = []

      lst_datasets = Dataset.all
      @snames = lst_datasets.pluck(:dataset_name)
      exps = experiments.to_a

      lst_datasets.each do |l|
         stdlst = exps.select {|e| e['dataset_id'] == l.dataset_id}
         @ccounts << stdlst.map{|c| c['cell_id']}.uniq.count
         @tcounts << stdlst.map {|c| c['tissue_id']}.uniq.count
      end

      if params[:c].to_s.strip.empty?
         @search_cell_lines = false
         @placeholder_cell_lines = 'Search cell line names ...'
      else
         @search_cell_lines = true
         if @cell_lines.length.equal? 0
            @caption_cell_lines = 'Your search for ' + params[:c].to_s.strip + " yielded no results."
            @placeholder_cell_lines = 'Search cell line names ...'
         else
            @caption_cell_lines = ''
            @placeholder_cell_lines = params[:c].to_s.strip
         end
      end

      # @gene_drug_associations = drug.gene_drugs.include(:gene, :dataset).order(:pvalue).limit(10)
      # sql = "SELECT s.*, g.*, gd.* FROM datasets s, genes g, gene_drugs gd WHERE gd.drug_id = #{@drug_id} AND s.dataset_id = gd.dataset_id AND g.gene_id = gd.gene_id ORDER BY gd.pvalue"
      # results = ActiveRecord::Base.connection.exec_query(sql)
      results = GeneDrug.includes(:gene, :dataset, :tissue).where(:drug_id => @drug_id).order(:pvalue)
      filtered_results = []
      results.each do |r|
        if params[:m].to_s.strip.empty?
          filtered_results << r
        else
          filtered_results << r if is_equal_incl(params[:m], r.gene['gene_name'])
        end
      end
      pseudotissues = ["pan-cancer", "solid tumour", "liquid tumour"]
      @gene_drug_associations = filtered_results.map{ |r| [r["dataset_id"], r.dataset["dataset_name"], r.gene["gene_name"], r["sens_stat"], r["estimate"], r["pvalue"], r.gene["gene_id"], r.tissue.nil? ? pseudotissues[-(r.tissue_id + 1)] : r.tissue.tissue_name, r["mDataType"]]}

      # rearrange gene drug associations
      if params[:sort] == "gene"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 2)
      elsif params[:sort] == "gd_dataset"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 1)
      elsif params[:sort] == "stat"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 3)
      elsif params[:sort] == "coef"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 4)
      elsif params[:sort] == "anova"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 5)
      else 
        # sort the experiments row by default for both tables
        @gene_drug_associations = sort_column(@gene_drug_associations, "asc", 5)
      end

      if params[:download] == "drug_table"
        data_csv = ''
        data_csv << 'mDataType,dataset_name,drug_name,sens_stat,estimate,pvalue,drug_id,tissue'
        data_csv << "\n"
        @gene_drug_associations.each do |d|
          data_csv << d[8].to_s + "," + d[1].to_s + "," + d[2].to_s + "," + d[3].to_s + "," + d[4].to_s + "," + d[5].to_s + "," + d[6].to_s + "," + d[7].to_s + "\n"
        end
        send_data data_csv, filename: @drug_name + '_gene_association_table.csv'
     end

      @gene_drug_associations = @gene_drug_associations.to_a.paginate(:page => params[:gdpage], :per_page => 10)

      if params[:download] == "cell_table"
         data_csv = ''
         data_csv << "Cell Name, Tissue, Datasets, # Experiments\n"
         cell_line_array.each do |d|
            data_csv << [0,3,1,2].map{|x| d[x]}.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @drug_name + '_cell_table.csv'
      end



      if params[:download] == "tissue_table"
         data_csv = ''
         data_csv << "Tissue Name, Datasets, # Experiments\n"
         tissue_array.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @drug_name + '_tissue_table.csv'
      end

      if params[:download] == "waterfall"
         data_csv = ''
         data_csv << "Cell, Dataset, AAC, IC50\n"
         waterfall_data.sort_by!{|g| g[0]}
         waterfall_data.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @drug_name + '_waterfall_table.csv'
      end


      @exit_code = 1
      return
   end #/drugs

   def genes
      @colspan = 2
      unless params[:id].present?
         genes = Gene.all
         @count_ids = []
         @count_names = []
         targets = Target.all
         targets.each do |t|
            sql = "SELECT DISTINCT dt.drug_id FROM targets t, drug_targets dt WHERE dt.target_id = #{t.target_id}"
            # genes = d.genes.pluck(:gene_id).uniq
            @count_ids << ActiveRecord::Base.connection.exec_query(sql).count
            @count_names << t.target_name
         end

         sql = "SELECT gene_id, gene_name FROM genes"
         genes = ActiveRecord::Base.connection.exec_query(sql)
         filtered_genes = []
         genes.each do |gene|
           if params[:g].to_s.strip.empty?
             filtered_genes << gene
           else
            filtered_genes << gene if is_equal_incl(params[:g], gene['gene_name'])
           end
         end
         @genes = filtered_genes.sort_by{ |g| g['gene_id']}.to_a.paginate(:page => params[:page], :per_page => 20)
         if @genes.length <= 10
            @colspan = 1
         end
         @genes_count = genes.count
         @exit_code = 0
         return
      end

      sql = "SELECT gene_name, ensg FROM genes WHERE gene_id = #{params[:id]}"
      gene = ActiveRecord::Base.connection.exec_query(sql)
      @gene_id = params[:id]
      @gene_name = gene[0]['gene_name']

      @synonyms = [gene[0]['ensg']]
      ## TODO:: Add protein based synonyms here!
      # sql = "SELECT drug_name, dataset_id, dataset_name FROM source_drug_names, datasets WHERE drug_id = #{@drug_id} AND dataset_id = source_id"
      # syn = ActiveRecord::Base.connection.exec_query(sql)
      # if syn.present?
      #    allsyn = []
      #    syn.each do |s|
      #       dataset_id = s['dataset_id']
      #       dataset_name = s['dataset_name']
      #       html_link = "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
      #       allsyn << [html_link, s['drug_name']]
      #    end
      #    names_so_far = []
      #    sources_so_far = []
      #    allsyn.each do |a|
      #       if names_so_far.include? a[1]
      #          if sources_so_far.include? a[0]
      #             next
      #          end
      #          @synonyms.each do |s|
      #             if s[1] == a[1]
      #                if s[0].include? a[0]
      #                   break
      #                else
      #                   s[0] = s[0] + ", " + a[0]
      #                   sources_so_far << a[0]
      #                   break
      #                end
      #             end
      #          end
      #       elsif sources_so_far.include? a[0]
      #          @synonyms.each do |s|
      #             if s[0] == a[0]
      #                s[1] = s[1] + ", " + a[1]
      #                names_so_far << a[1]
      #                break
      #             end
      #          end
      #          unless names_so_far.include? a[1]
      #             @synonyms << a
      #             sources_so_far << a[0]
      #             names_so_far << a[1]
      #          end
      #       else
      #          @synonyms << a
      #          sources_so_far << a[0]
      #          names_so_far << a[1]
      #       end
      #    end
      # end

      unless @synonyms.present?
         @synonyms = "N/A"
      end

      array = []
      tissue_array = []
      drugs = []
      tissues = []
      datasets = []

      sql = "SELECT d.drug_id, d.drug_name, s.dataset_id, s.dataset_name FROM drugs d, datasets s, drug_targets dt, experiments e, targets t WHERE dt.target_id = t.target_id AND t.gene_id = #{@gene_id} AND s.dataset_id = e.dataset_id AND dt.drug_id = e.drug_id AND d.drug_id = dt.drug_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)

      if experiments.present?
         experiments.each do |e|
            drug_id = e['drug_id']
            drug = e['drug_name']
            dataset_id = e['dataset_id']
            dataset = e['dataset_name']
            if drugs.include? drug
               array.each do |a|
                  if a[0] == drug
                     unless a[1].include? dataset
                        a[1] = a[1] + ", " + "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>"
                        datasets << dataset
                     end
                     a[2] += 1
                     break
                  end
               end
            else
              if params[:d].to_s.strip.empty?
                array << [drug, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, drug_id]
              else
                array << [drug, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, drug_id] if is_equal_incl(params[:d], drug)
              end
              drugs << drug
              datasets << dataset
            end
            # if tissues.include? tissue
            #    tissue_array.each do |t|
            #       if t[0] == tissue
            #          unless t[1].include? dataset
            #             t[1] = t[1] + ", " + "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>"
            #             datasets << dataset
            #          end
            #          t[2] += 1
            #          break
            #       end
            #    end
            # else
            #   if params[:t].to_s.strip.empty?
            #     tissue_array << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id]
            #   else
            #     tissue_array << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id] if is_equal_incl(params[:t], tissue)
            #   end
            #   tissues << tissue
            #   datasets << dataset
            # end
         end
      end

      if params[:sort] == "compounds"
        array = sort_column(array, params[:direction], 0)
      elsif params[:sort] == "datasets"
        array = sort_column(array, params[:direction], 1)
      elsif params[:sort] == "experiments"
        array = sort_column(array, params[:direction], 2)
      else
        # sort the experiments row by default for both tables
        array = sort_column(array, "desc", 2)
      end

      @num_drugs = drugs.uniq.count
      # @num_tissues = tissues.uniq.count
      @numdatasets = datasets.uniq.count

      # tissue_array = tissue_array.sort_by { |e| e[2] }.reverse
      @drugs = array.paginate(:page => params[:dpage], :per_page => 10)
      # @tissues = tissue_array.paginate(:page => params[:tpage], :per_page => 10)

      # sql = "SELECT s.*, d.*, gd.* FROM datasets s, drugs d, gene_drugs gd WHERE gd.gene_id = #{@gene_id} AND s.dataset_id = gd.dataset_id AND d.drug_id = gd.drug_id ORDER BY gd.pvalue"
      # results = ActiveRecord::Base.connection.exec_query(sql)
      results = GeneDrug.includes(:drug, :dataset, :tissue).where(:gene_id => @gene_id).order(:pvalue)
      filtered_results = []
      results.each do |r|
        if params[:m].to_s.strip.empty?
          filtered_results << r
        else
          filtered_results << r if is_equal_incl(params[:m], r.drug['drug_name'])
        end
      end
      pseudotissues = ["pan-cancer", "solid tumour", "liquid tumour"]
      @gene_drug_associations = filtered_results.map{ |r| [r["dataset_id"], r.dataset["dataset_name"], r.drug["drug_name"], r["sens_stat"], r["estimate"], r["pvalue"], r["drug_id"], r.tissue.nil? ? pseudotissues[-(r.tissue_id + 1)] : r.tissue.tissue_name, r["mDataType"]]}

      # rearrange gene drug associations
      if params[:sort] == "drug"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 2)
      elsif params[:sort] == "gd_dataset" 
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 1)
      elsif params[:sort] == "stat"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 3)
      elsif params[:sort] == "coef"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 4)
      elsif params[:sort] == "anova"
        @gene_drug_associations = sort_column(@gene_drug_associations, params[:direction], 5)
      else 
        # sort the experiments row by default for both tables
        @gene_drug_associations = sort_column(@gene_drug_associations, "asc", 5)
      end

      if params[:download] == "drug_table"
        data_csv = ''
        data_csv << 'mDataType,dataset_name,drug_name,sens_stat,estimate,pvalue,drug_id,tissue'
        data_csv << "\n"
        @gene_drug_associations.each do |d|
           data_csv << d[8].to_s + "," + d[1].to_s + "," + d[2].to_s + "," + d[3].to_s + "," + d[4].to_s + "," + d[5].to_s + "," + d[6].to_s + "," + d[7].to_s + "\n"
        end
        send_data data_csv, filename: @gene_name + '_drug_association_table.csv'
     end

      @gene_drug_associations = @gene_drug_associations.to_a.paginate(:page => params[:gdpage], :per_page => 10)

      @ccounts = []
      # @tcounts = []

      lst_datasets = Dataset.all
      @snames = lst_datasets.pluck(:dataset_name)
      exps = experiments.to_a

      lst_datasets.each do |l|
         stdlst = exps.select {|e| e['dataset_id'] == l.dataset_id}
         @ccounts << stdlst.map{|d| d['drug_id']}.uniq.count
         # @tcounts << stdlst.map {|c| c['tissue_id']}.uniq.count
      end

      if params[:d].to_s.strip.empty?
         @search_drugs = false
         @placeholder_drug = 'Search compound names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drugs = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search compound names ...'
         else
            @caption_drugs = ''
            @placeholder_drug = params[:d].to_s.strip
         end
      end

      @exit_code = 1
      return
   end #/gene

   def datasets
      unless params[:id].present?
         sql = "SELECT dataset_id, dataset_name FROM datasets"
         datasets = ActiveRecord::Base.connection.exec_query(sql)
         @datasets = datasets.to_a.paginate(:page => params[:page], :per_page => 20)
         @datasets_count = datasets.count
         @exit_code = 0
         return
      end

      sql = "SELECT dataset_name FROM datasets WHERE dataset_id = #{params[:id]}"
      dataset = ActiveRecord::Base.connection.exec_query(sql)
      @dataset_id = params[:id]
      @dataset_name = dataset[0]['dataset_name']

      @path = 'datasets/' + @dataset_name.downcase

      sql = "SELECT c.cell_id, c.cell_name, d.drug_id, d.drug_name FROM cells c, drugs d, experiments e WHERE e.dataset_id = #{@dataset_id} AND c.cell_id = e.cell_id AND d.drug_id = e.drug_id"
      cds = ActiveRecord::Base.connection.exec_query(sql)
      cell_lines = []
      drugs = []

      cds.each do |row|
        if params[:c].nil? || params[:c].to_s.strip.empty? # using lazy evaluation to remove need for to_s on nil. 
          cell_lines << [row['cell_id'], row['cell_name']]
        else
          cell_lines << [row['cell_id'], row['cell_name']] if is_equal_incl(params[:c], row['cell_name'].to_s)
        end
        if params[:d].nil? || params[:d].to_s.strip.empty?
          drugs << [row['drug_id'], row['drug_name']]
        else
          drugs << [row['drug_id'], row['drug_name']] if is_equal_incl(params[:d], row['drug_name'].to_s)
        end
      end
      cell_lines = cell_lines.uniq.sort
      drugs = drugs.uniq.sort
      #drugs = cds.map {|c| [c['drug_id'], c['drug_name']]}.uniq.sort
      @cell_lines_count = cell_lines.count
      @drugs_count = drugs.count
      @cell_lines = cell_lines.paginate(:page => params[:cpage], :per_page => 20)
      @drugs = drugs.paginate(:page => params[:dpage], :per_page => 20)

      lst_datasets = Dataset.all
      @snames = lst_datasets.pluck(:dataset_name)
      @ccounts = []
      @tcounts = []
      @dcounts = []
      @ecounts = []
      @colors = []

      lst_datasets.each do |l|
        #  sql = "SELECT e.cell_id, e.drug_id, t.tissue_id FROM tissues t, experiments e, cell_tissues ct WHERE e.dataset_id = #{l.dataset_id} AND ct.cell_id = e.cell_id AND t.tissue_id = ct.tissue_id"
        #  exp = ActiveRecord::Base.connection.exec_query(sql)
         # @ccounts << exp.map { |e| e['cell_id']  }.uniq.count
         # @tcounts << exp.map { |e| e['tissue_id']  }.uniq.count
        #  @dcounts << exp.map { |e| e['drug_id']  }.uniq.count
        @dcounts << Dataset.joins(:drugs).where(:dataset_id=>l.dataset_id).group("drugs.drug_id").length
        @ecounts << Experiment.where(:dataset_id => l.dataset_id).count
        # exps  = Experiment.where(:dataset_id => l.dataset_id).pluck(:cell_id, :tissue_id).uniq
        @ccounts << Dataset.joins(:cells).where(:dataset_id=>l.dataset_id).count
        @tcounts << Dataset.joins(:tissues).where(:dataset_id=>l.dataset_id).group("tissues.tissue_id").length
         if l.dataset_id.to_s == @dataset_id.to_s
            @colors << 'rgba(222,45,38,0.8)'
         else
            @colors << 'rgb(142,124,195)'
         end
      end


      # lst_datasets.each do |l|
      #    sql = "SELECT dc.cell_id, ct.tissue_id FROM cell_tissues ct, dataset_cells dc WHERE dc.dataset_id = #{l.dataset_id} AND ct.cell_id = dc.cell_id"
      #    exp = ActiveRecord::Base.connection.exec_query(sql)
      #    @ccounts << exp.map { |e| e['cell_id']  }.uniq.count
      #    @tcounts << exp.map { |e| e['tissue_id']  }.uniq.count
      #    # @dcounts << exp.map { |e| e['drug_id']  }.uniq.count
      #    # @ecounts << exp.count
      #    if l.dataset_id.to_s == @dataset_id.to_s
      #       @colors << 'rgba(222,45,38,0.8)'
      #    else
      #       @colors << 'rgb(142,124,195)'
      #    end
      # end




      # Set right flag for table caption

      if params[:d].to_s.strip.empty?
         @search_drugs = false
         @placeholder_drug = 'Search compound names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drug = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search compound names ...'
         else
            @caption_drug = ''
            @placeholder_drug = params[:d].to_s.strip
         end
      end


      if params[:c].to_s.strip.empty?
         @search_cell_lines = false
         @placeholder_cell_lines = 'Search cell line names ...'
      else
         @search_cell_lines = true
         if @cell_lines.length.equal? 0
            @caption_cell_lines = 'Your search for ' + params[:c].to_s.strip + " yielded no results."
            @placeholder_cell_lines = 'Search cell line names ...'
         else
            @caption_cell_lines = ''
            @placeholder_cell_lines = params[:c].to_s.strip
         end
      end



      if params[:download] == "cell_table"
         data_csv = ''
         data_csv << "Cell Name\n"
         cell_lines.each do |d|
            data_csv << d[1]
            data_csv << "\n"
         end
         send_data data_csv, filename: @dataset_name + '_cell_table.csv'
      end

      if params[:download] == "drug_table"
         data_csv = ''
         data_csv << "Drug Name\n"
         drugs.each do |d|
            data_csv << d[1]
            data_csv << "\n"
         end
         send_data data_csv, filename: @dataset_name + '_drug_table.csv'
      end



      @exit_code = 1
      return
   end #/datasets

   def experiments
     datasets = Dataset.all
     @dataset_names = []
     @exp_per_cell_line = []
     @exp_per_drug = []
     datasets.each do |d|
        sql = "SELECT cell_id, COUNT(DISTINCT experiment_id) FROM experiments WHERE dataset_id = #{d.dataset_id} GROUP BY cell_id"
        result = ActiveRecord::Base.connection.exec_query(sql)
        sum = 0
        result.each { |r| sum += r['COUNT(DISTINCT experiment_id)'] }
        sum = sum / result.count
        @exp_per_cell_line << sum

        sql = "SELECT drug_id, COUNT(DISTINCT experiment_id) FROM experiments WHERE dataset_id = #{d.dataset_id} GROUP BY drug_id"
        result = ActiveRecord::Base.connection.exec_query(sql)
        sum = 0
        result.each { |r| sum += r['COUNT(DISTINCT experiment_id)'] }
        sum = sum / result.count
        @exp_per_drug << sum

        @dataset_names << d.dataset_name
     end
   end #/experiments

   def explore
    
    @tissues = ActiveRecord::Base.connection.exec_query("SELECT tissue_id, tissue_name FROM tissues")
    @cell_lines = ActiveRecord::Base.connection.exec_query("SELECT cell_id, cell_name FROM cells")
    @drugs = ActiveRecord::Base.connection.exec_query("SELECT drug_id, drug_name FROM drugs")
    @drug_targets = ActiveRecord::Base.connection.exec_query("SELECT target_id, target_name FROM targets")

    @valid_cell_line_ids = []
    @valid_drug_ids = []

    @is_tissues = true
    @is_cell_lines = false
    @is_drugs = false    
    @is_drug_targets = true

    if params[:select_all].present?
      val = params[:select_all].to_s.split(' ')
      if val[0] == "true"
        @is_cell_lines = true
        @is_drugs = true
      end
    end

    if params[:tid].present?  # tissues
      tids = params[:tid].to_s.split(' ')
      @tissues = @tissues.select{|t| tids.include? t['tissue_id'].to_s}
      @is_cell_lines = true
      @is_drugs = true
      # cell lines
      sql = "SELECT DISTINCT cell_id, cell_name FROM cells WHERE tissue_id IN (#{tids.join(', ')})"
      @cell_lines = ActiveRecord::Base.connection.exec_query(sql)
      #drugs
      sql = "SELECT DISTINCT drugs.drug_id, drugs.drug_name FROM experiments INNER JOIN drugs WHERE experiments.tissue_id IN (#{tids.join(', ')}) AND experiments.drug_id = drugs.drug_id"
      @drugs = ActiveRecord::Base.connection.exec_query(sql)
    end

    if params[:cid].present? # cell lines
      cids = params[:cid].to_s.split(' ')
      sql = "SELECT DISTINCT drug_id FROM experiments WHERE cell_id IN (#{cids.join(', ')})"
      drug_ids = []
      results = ActiveRecord::Base.connection.exec_query(sql)
      results.to_hash.each{ |k, v| drug_ids.push(k.values[0]) }
      respond_to do |format|
        response = { :drug_ids => drug_ids }
        format.json  { render :json => response }
      end
    end

    if params[:drug_id].present? # drugs
      drug_ids = params[:drug_id].to_s.split(' ')
      sql = "SELECT DISTINCT cell_id FROM experiments WHERE drug_id IN (#{drug_ids.join(', ')})"
      cids = []
      results = ActiveRecord::Base.connection.exec_query(sql)
      results.to_hash.each{ |k, v| cids.push(k.values[0]) }
      respond_to do |format|
        response = { :cids => cids }
        format.json  { render :json => response }
      end

    end

    if params[:target_id].present?  # targets
      target_ids = params[:target_id].to_s.split(' ')
      @drug_targets = @drug_targets.select{|t| target_ids.include? t['target_id'].to_s}
      @is_cell_lines = true
      @is_drugs = true
      # cell lines
      sql = "SELECT DISTINCT c.cell_id, c.cell_name FROM drug_targets d, targets t, experiments e, cells c WHERE t.target_id IN (#{target_ids.join(', ')}) AND e.cell_id = c.cell_id AND e.drug_id = d.drug_id AND d.target_id = t.target_id"
      @cell_lines = ActiveRecord::Base.connection.exec_query(sql)
      # drugs
      sql = "SELECT DISTINCT dr.drug_id, dr.drug_name FROM drug_targets d, targets t, experiments e, drugs dr WHERE t.target_id IN (#{target_ids.join(', ')}) AND e.drug_id = dr.drug_id AND e.drug_id = d.drug_id AND d.target_id = t.target_id"
      @drugs = ActiveRecord::Base.connection.exec_query(sql)
    end

    @cell_lines_count = @cell_lines.count
    @tissues_count = @tissues.count
    @drugs_count = @drugs.count
    @drug_targets_count = @drug_targets.count
   end

   def contact_us
     @title = params[:title].present? ? params[:title].to_s.strip : ""
     @subject = params[:subject].present? ? params[:subject].to_s.strip : ""
     @id = params[:id].present? ? params[:id] : "-1"
   end

end
