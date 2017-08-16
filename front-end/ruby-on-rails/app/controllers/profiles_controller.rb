class ProfilesController < ApplicationController

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
            html_link = "<a target=\"_blank\" href=\"http://web.expasy.org/cellosaurus/#{accession}\">Cellosaurus</a>"
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

      sql = "SELECT d.drug_id, d.drug_name, s.dataset_id, s.dataset_name, e.experiment_id FROM drugs d, datasets s, experiments e WHERE e.cell_id = #{@cell_line_id} AND d.drug_id = e.drug_id AND s.dataset_id = e.dataset_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)

      drugs = []
      datasets = []
      dataset_count = []
      drugs_arr = []

      @drug_AAC = []
      @drug_IC50 = []

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
            @drug_AAC[ind] << Experiment.find(e['experiment_id'] ).profile["AAC"]
            if Experiment.find(e['experiment_id'] ).profile["IC50"].nil?
              temp = Float::NAN
            else
              temp = Experiment.find(e['experiment_id'] ).profile["IC50"]
            end
            @drug_IC50[ind] << temp
         else
            drugs << drug_name
            datasets << dataset_name
            @drug_AAC << [Experiment.find(e['experiment_id'] ).profile["AAC"]]

            if Experiment.find(e['experiment_id'] ).profile["IC50"].nil?
              temp = Float::NAN
            else
              temp = Experiment.find(e['experiment_id'] ).profile["IC50"]
            end
            @drug_IC50 << [temp]
            if params[:d].to_s.strip.empty?
              drugs_arr << [drug_name, "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>", 1, drug_id]
            else
              drugs_arr << [drug_name, "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>", 1, drug_id] if is_equal_incl(params[:d], drug_name)
            end
         end
      end

      @drug_names_waterfall = drugs
      @numdatasets = datasets.uniq.count
      @numdrugs = drugs.count
      drugs_arr = drugs_arr.sort_by { |e| e[2] }.reverse 

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
         @placeholder_drug = 'Search drug names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drug = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search drug names ...'
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

      sql = "SELECT tissue_name FROM tissues WHERE tissue_id = #{params[:id]}"
      tissue = ActiveRecord::Base.connection.exec_query(sql)
      @tissue_id = params[:id]
      @tissue_name = tissue[0]['tissue_name']

      @synonyms = []

      sql = "SELECT tissue_name, dataset_id, dataset_name FROM source_tissue_names, datasets WHERE tissue_id = #{@tissue_id} AND dataset_id = source_id"
      syn = ActiveRecord::Base.connection.exec_query(sql)
      if syn.present?
         allsyn = []
         syn.each do |s|
            dataset_id = s['dataset_id']
            dataset_name = s['dataset_name']
            html_link = "<a href=\"/datasets/#{dataset_id}\">" + dataset_name + "</a>"
            allsyn << [html_link, s['tissue_name']]
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

      @numdrugs = drugs.uniq.count
      @numdatasets = datasets.uniq.count
      array = array.sort_by { |e| e[2] }.reverse
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


      sql = "SELECT ct.cell_id, s.dataset_id FROM datasets s, experiments e, cell_tissues ct WHERE ct.tissue_id = #{@tissue_id} AND e.cell_id = ct.cell_id AND s.dataset_id = e.dataset_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)


      exps = experiments.to_a
      @ccounts = []
      datasets.each do |s|
         count = exps.select {|e| e['dataset_id'] == s['dataset_id']}.map{|x| x['cell_id']}.uniq.count
         cells_count << [s['dataset_name'],count]
      end

      @ccounts = cells_count.map {|x| x[1] }

      # Set right flag for table caption

      if params[:d].to_s.strip.empty?
         @search_drugs = false
         @placeholder_drug = 'Search drug names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drug = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search drug names ...'
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
         # XXX :: The database was changes as a workaround to get nice drugs first. Please fix the actual data!
         @drugs = drugs.sort_by{ |d| d['drug_id']}.to_a.paginate(:page => params[:page], :per_page => 20)
         @drugs_count = drugs.count
         @exit_code = 0
         return
      end

      sql = "SELECT drug_name FROM drugs WHERE drug_id = #{params[:id]}"
      drug = ActiveRecord::Base.connection.exec_query(sql)
      @drug_id = params[:id]
      @drug_name = drug[0]['drug_name']

      @synonyms = []
      @drug_ids = DrugAnnot.where(:drug_id => @drug_id).first
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

      array = []
      tarray = []
      cell_lines = []
      tissues = []
      datasets = []

      @cell_AAC = []
      @cell_IC50 = []

      sql = "SELECT c.cell_id, c.cell_name, t.tissue_id, t.tissue_name, s.dataset_id, s.dataset_name, e.experiment_id FROM cells c, tissues t, datasets s, cell_tissues ct, experiments e WHERE e.drug_id = #{@drug_id} AND s.dataset_id = e.dataset_id AND ct.cell_id = e.cell_id AND t.tissue_id = ct.tissue_id AND c.cell_id = e.cell_id"
      experiments = ActiveRecord::Base.connection.exec_query(sql)

      if experiments.present?
         experiments.each do |e|
            cell_id = e['cell_id']
            cell_line = e['cell_name']
            tissue_id = e['tissue_id']
            tissue = e['tissue_name']
            dataset_id = e['dataset_id']
            dataset = e['dataset_name']
            if cell_lines.include? cell_line
               array.each do |a|
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
               @cell_AAC[ind] << Experiment.find(e['experiment_id'] ).profile["AAC"]
               if Experiment.find(e['experiment_id'] ).profile["IC50"].nil?
                 temp = Float::NAN
               else
                 temp = Experiment.find(e['experiment_id'] ).profile["IC50"]
               end
               @cell_IC50[ind] << temp
            else
              @cell_AAC << [Experiment.find(e['experiment_id'] ).profile["AAC"]]

              if Experiment.find(e['experiment_id'] ).profile["IC50"].nil?
                temp = Float::NAN
              else
                temp = Experiment.find(e['experiment_id'] ).profile["IC50"]
              end
              @cell_IC50 << [temp]
              if params[:c].to_s.strip.empty?
                array << [cell_line, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, cell_id, tissue]
              else
                array << [cell_line, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, cell_id, tissue] if is_equal_incl(params[:c], cell_line)
              end
              cell_lines << cell_line
              datasets << dataset
            end
            if tissues.include? tissue
               tarray.each do |t|
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
                tarray << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id]
              else
                tarray << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id] if is_equal_incl(params[:t], tissue)
              end
              tissues << tissue
              datasets << dataset
            end
         end
      end

      @cell_lines_waterfall = cell_lines
      @num_cell_lines = cell_lines.uniq.count
      @num_tissues = tissues.uniq.count
      @numdatasets = datasets.uniq.count
      array = array.sort_by { |e| e[2] }.reverse
      tarray = tarray.sort_by { |e| e[2] }.reverse
      @cell_lines = array.paginate(:page => params[:cpage], :per_page => 10)
      @tissues = tarray.paginate(:page => params[:tpage], :per_page => 10)



      if params[:download] == "cell_table"
         array.each do |d| 
            d[1] = d[1].gsub(/<a[^>]*>/, '')
            d[1] = d[1].gsub(/<\/a.?>/, '')
            d[1] = d[1].gsub(/,/, ';')
            d.delete_at(3)
         end
      end

      if params[:download] == "tissue_table"
         tarray.each do |d| 
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



      if params[:download] == "cell_table"
         data_csv = ''
         data_csv << "Cell Name, Tissue, Datasets, # Experiments\n"
         array.each do |d|
            data_csv << [0,3,1,2].map{|x| d[x]}.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @drug_name + '_cell_table.csv'
      end



      if params[:download] == "tissue_table"
         data_csv = ''
         data_csv << "Tissue Name, Datasets, # Experiments\n"
         tarray.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @drug_name + '_tissue_table.csv'
      end


      @exit_code = 1
      return
   end #/drugs

   def targets
      unless params[:id].present?
         targets = Target.all
         @count_ids = []
         @count_names = []
         targets.each do |t|
            sql = "SELECT DISTINCT dt.drug_id FROM targets t, drug_targets dt WHERE dt.target_id = #{t.target_id}"
            # targets = d.targets.pluck(:target_id).uniq
            @count_ids << ActiveRecord::Base.connection.exec_query(sql).count
            @count_names << t.target_name
         end

         sql = "SELECT target_id, target_name FROM targets"
         targets = ActiveRecord::Base.connection.exec_query(sql)
         @targets = targets.sort_by{ |d| d['drug_id']}.to_a.paginate(:page => params[:page], :per_page => 20)
         @targets_count = targets.count
         @exit_code = 0
         return
      end

      sql = "SELECT target_name FROM targets WHERE target_id = #{params[:id]}"
      target = ActiveRecord::Base.connection.exec_query(sql)
      @target_id = params[:id]
      @target_name = target[0]['target_name']

      @synonyms = []
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
      tarray = []
      drugs = []
      tissues = []
      datasets = []

      sql = "SELECT d.drug_id, d.drug_name, s.dataset_id, s.dataset_name FROM drugs d, datasets s, drug_targets dt, experiments e WHERE dt.target_id = #{@target_id} AND s.dataset_id = e.dataset_id AND dt.drug_id = e.drug_id AND d.drug_id = dt.drug_id"
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
            #    tarray.each do |t|
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
            #     tarray << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id]
            #   else
            #     tarray << [tissue, "<a href=\"/datasets/#{dataset_id}\">" + dataset + "</a>", 1, tissue_id] if is_equal_incl(params[:t], tissue)
            #   end
            #   tissues << tissue
            #   datasets << dataset
            # end
         end
      end

      @num_drugs = drugs.uniq.count
      # @num_tissues = tissues.uniq.count
      @numdatasets = datasets.uniq.count
      array = array.sort_by { |e| e[2] }.reverse
      # tarray = tarray.sort_by { |e| e[2] }.reverse
      @drugs = array.paginate(:page => params[:dpage], :per_page => 10)
      # @tissues = tarray.paginate(:page => params[:tpage], :per_page => 10)

      if params[:download] == "drug_table"
         array.each do |d| 
            d[1] = d[1].gsub(/<a[^>]*>/, '')
            d[1] = d[1].gsub(/<\/a.?>/, '')
            d[1] = d[1].gsub(/,/, ';')
            d.delete_at(3)
         end
      end


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
         @placeholder_drug = 'Search drug names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drugs = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search drug names ...'
         else
            @caption_drugs = ''
            @placeholder_drug = params[:d].to_s.strip
         end
      end


      if params[:download] == "drug_table"
         data_csv = ''
         data_csv << "Drug Name, Datasets, # Experiments\n"
         array.each do |d|
            data_csv << d.to_csv
            data_csv << "\n"
         end
         send_data data_csv, filename: @target_name + '_drug_table.csv'
      end


      @exit_code = 1
      return
   end #/target

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
      cds.each do |row|
        if params[:c].to_s.strip.empty?
          cell_lines << [row['cell_id'], row['cell_name']]
        else
          cell_lines << [row['cell_id'], row['cell_name']] if is_equal_incl(params[:c], row['cell_name'].to_s)
        end
      end
      cell_lines = cell_lines.uniq.sort
      drugs = []
      cds.each do |row|
        if params[:d].to_s.strip.empty?
          drugs << [row['drug_id'], row['drug_name']]
        else
          drugs << [row['drug_id'], row['drug_name']] if is_equal_incl(params[:d], row['drug_name'].to_s)
        end
      end
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
         sql = "SELECT e.cell_id, e.drug_id, t.tissue_id FROM tissues t, experiments e, cell_tissues ct WHERE e.dataset_id = #{l.dataset_id} AND ct.cell_id = e.cell_id AND t.tissue_id = ct.tissue_id"
         exp = ActiveRecord::Base.connection.exec_query(sql)
         @ccounts << exp.map { |e| e['cell_id']  }.uniq.count
         @tcounts << exp.map { |e| e['tissue_id']  }.uniq.count
         @dcounts << exp.map { |e| e['drug_id']  }.uniq.count
         @ecounts << exp.count
         if l.dataset_id.to_s == @dataset_id.to_s
            @colors << 'rgba(222,45,38,0.8)'
         else
            @colors << 'rgb(142,124,195)'
         end
      end


      # Set right flag for table caption

      if params[:d].to_s.strip.empty?
         @search_drugs = false
         @placeholder_drug = 'Search drug names ...'
      else
         @search_drugs = true
         if @drugs.length.equal? 0
            @caption_drug = 'Your search for ' + params[:d].to_s.strip + " yielded no results."
            @placeholder_drug = 'Search drug names ...'
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
    @is_cell_lines = false
    @drugs = ActiveRecord::Base.connection.exec_query("SELECT drug_id, drug_name FROM drugs")
    @is_drugs = false
    @drug_targets = ActiveRecord::Base.connection.exec_query("SELECT target_id, target_name FROM targets")
    @is_tissues = true
    @is_drug_targets = true
    if params[:tid].present?  # tissues
      tids = params[:tid].to_s.split('+')
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
    if params[:target_id].present?  # targets
      target_ids = params[:target_id].to_s.split('+')
      @drug_targets = @drug_targets.select{|t| target_ids.include? t['target_id'].to_s}
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
     @title = params[:title].present? ? params[:title] : ""
     @subject = params[:subject].present? ? params[:subject] : ""

     @title = @title.to_s.strip
     @subject = @subject.to_s.strip

   end

end
