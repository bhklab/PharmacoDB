class AutocompleteController < ApplicationController

    ### FIXME:: navitoclax:birinapant (1:1 m 
    ## returns no results.
    ## This whole thing is basically broken, does not work on drugs with special characters, 
    ## terms with more than one space in them, etc... 

    def suggest
        # puts params.values[0]
        input = params.values[0].split(Rails.application.config.my_search_split).last
        previous = params.values[0].split(Rails.application.config.my_search_split)[0...-1]
        # input = params.values[0].split(" ").last
        # previous = params.values[0].split(" ")[0...-1]
        # background = params.values[0].split(" ")[0...-1].join(" ")
        background = ""
        previous_types = []
        previous_ids = []
        allowed_searches = []
        # puts params.values[0].split(Rails.application.config.my_search_split)
        # puts params.values[0].split(" ")
        # puts input
        # puts previous
        if previous.empty?
            allowed_searches << "all"
        else
            term = previous.shift
            matched = false
            pc = previous.clone
            while !term.nil?
                cells1 = Cell.where(:cell_name => term).pluck(:cell_id)
                cells2 = SourceCellName.where(:cell_name => term).pluck(:cell_id)
                cells = cells1 + cells2

                if !cells.empty?
                    previous_types << 'cell'
                    previous_ids << cells.uniq
                    matched = true
                    # next
                end

                drugs1 = Drug.where(:drug_name => term).pluck(:drug_id)
                drugs2 = SourceDrugName.where(:drug_name => term).pluck(:drug_id)
                drugs = drugs1 + drugs2


                if !drugs.empty?
                    previous_types << 'drug'
                    previous_ids << drugs.uniq
                    matched = true
                    # next
                end

                targets = Target.where(:target_name => term).pluck(:target_id)


                if !targets.empty?
                    previous_types << 'target'
                    previous_ids << targets.uniq
                    matched = true
                    # next
                end


                tissues1 = Tissue.where(:tissue_name => term).pluck(:tissue_id)
                tissues2 = SourceTissueName.where(:tissue_name => term).pluck(:tissue_id)
                tissues = tissues1 + tissues2


                if !tissues.empty?
                    previous_types << 'tissue'
                    previous_ids << tissues.uniq
                    matched = true
                    # next
                end

                datasets = Dataset.where(:dataset_name => term).pluck(:dataset_id)
            
                if !datasets.empty?
                    previous_types << 'dataset'
                    previous_ids << datasets.uniq
                    matched = true
                    # next
                end

                if matched
                    # puts term
                    background += term + " "
                    matched = false
                    previous = pc
                    term = previous.shift
                    pc = previous.clone
                else
                    next_term = pc.shift
                    if next_term.nil?
                        input = term + " " + input
                        term = previous.shift
                    else
                        term += " " + next_term
                    end
                end


            end

            if previous_types.empty?
                allowed_searches << "all"
            # elsif previous_types.length != previous.length
            #     allowed_searches << "none"
            elsif previous_types.to_set == ["dataset"].to_set
                allowed_searches << "dataset"
                if previous_types.length == 1
                    allowed_searches << "drug"
                    allowed_searches << "cell"
                end
            elsif previous_types == ["drug"]
                allowed_searches << "cell"
            elsif previous_types == ["cell"]
                allowed_searches << "drug"
            elsif previous_types.to_set == ["cell", "drug"].to_set
                allowed_searches << "dataset"
            end
        end
        # puts input
        # puts previous_types


        # puts allowed_searches 
        


        if params[:cell_query].present?

            unless allowed_searches.include?('all') || allowed_searches.include?('cell')
                @cells = []
                render json: @cells.map{|i| i.downcase}.uniq.to_json
                return
            end
            
            ## TEST CASE: 253J + 17-AAG
            cells1 = Cell.where("cell_name LIKE ?", "#{input}%")
            cells2 = SourceCellName.where("source_cell_names.cell_name LIKE ?", "#{input}%")
            # cells2 = SourceCellName.where("cell_name LIKE ?", "#{input}%")
            if (previous_types.include? "drug")
                drug_ids = previous_ids[previous_types.find_index{|x| x=="drug"}]
                cells1 = cells1.includes(:experiments).where(:experiments => {:drug_id => drug_ids})
                cells2 = cells2.includes(:experiments).where(:experiments => {:drug_id => drug_ids})
            end
          #   if (previous_types.include? "dataset")
                # dataset_ids = previous_ids[previous_types.find_index{|x| x=="dataset"}]
             #    cells1 = cells1.includes(:experiments).where(:experiments => {:dataset_id => dataset_ids})
             #    cells2 = cells2.includes(:experiments).where(:experiments => {:dataset_id => dataset_ids})
          #   end


            cells1 = cells1.pluck(:cell_name)
            cells2 = cells2.pluck(:cell_name)

            cells = cells1 + cells2
            unless background == ""
                @cells = cells.map{|x| background + x}
            end
            if background == ""
                @cells = cells
            end
            render json: @cells.map{|i| i.downcase}.uniq.to_json
            return
        elsif params[:tissue_query].present?

            unless allowed_searches.include?('all') || allowed_searches.include?('tissue')
                @tissues = []
                render json: @tissues.map{|i| i.downcase}.uniq.to_json
                return
            end

            tissues1 = Tissue.where("tissue_name LIKE ?", "#{input}%").pluck(:tissue_name)
            tissues2 = SourceTissueName.where("tissue_name LIKE ?", "#{input}%").pluck(:tissue_name)
            tissues = tissues1 + tissues2
            unless background == ""
                @tissues = tissues.map{|x| background + x}
            end
            if background == ""
                @tissues = tissues
            end
            render json: @tissues.map{|i| i.downcase}.uniq.to_json
            return
        elsif params[:drug_query].present?

            unless allowed_searches.include?('all') || allowed_searches.include?('drug')
                @drugs = []
                render json: @drugs.map{|i| i.downcase}.uniq.to_json
                return
            end

            drugs1 = Drug.where("drug_name LIKE ?", "#{input}%")
            drugs2 = SourceDrugName.where("source_drug_names.drug_name LIKE ?", "#{input}%")

            if (previous_types.include? "cell")
                cell_ids = previous_ids[previous_types.find_index{|x| x=="cell"}]
                drugs1 = drugs1.includes(:experiments).where(:experiments => {:cell_id => cell_ids})
                drugs2 = drugs2.includes(:experiments).where(:experiments => {:cell_id => cell_ids})
            end
            drugs1 = drugs1.pluck(:drug_name)
            drugs2 = drugs2.pluck(:drug_name)

            drugs = drugs1 + drugs2
            unless background == ""
                @drugs = drugs.map{|x| background + x}
            end
            if background == ""
                @drugs = drugs
            end
            render json: @drugs.map{|i| i.downcase}.uniq.to_json
            return
        elsif params[:target_query].present?

            unless allowed_searches.include?('all') || allowed_searches.include?('target')
                @targets = []
                render json: @targets.map{|i| i.downcase}.uniq.to_json
                return
            end

            targets = Target.where("target_name LIKE ?", "#{input}%")

            targets = targets.pluck(:target_name)

            unless background == ""
                @targets = targets.map{|x| background + x}
            end
            if background == ""
                @targets = targets
            end
            render json: @targets.uniq.to_json
            return
            
        elsif params[:dataset_query].present?

            unless allowed_searches.include?('all') || allowed_searches.include?('dataset')
                @datasets = []
                render json: @datasets.map{|i| i.downcase}.uniq.to_json
                return
            end

            datasets = Dataset.where("dataset_name LIKE ?", "#{input}%")

            if (previous_types.include? "cell")
                cell_ids = previous_ids[previous_types.find_index{|x| x=="cell"}]
                datasets = datasets.includes(:experiments).where(:experiments => {:cell_id => cell_ids})
            end

            if (previous_types.include? "drug")
                drug_ids = previous_ids[previous_types.find_index{|x| x=="drug"}]
                datasets = datasets.includes(:experiments).where(:experiments => {:drug_id => drug_ids})
            end

            datasets = datasets.pluck(:dataset_name)

            unless background == ""
                @datasets = datasets.map{|x| background + x}
            end
            if background == ""
                @datasets = datasets
            end
            render json: @datasets.uniq.to_json
            return
        elsif params[:other_query].present?
            cells1 = Cell.where("cell_name LIKE ?", "#{input}%").pluck(:cell_name)
            cells2 = SourceCellName.where("cell_name LIKE ?", "#{input}%").pluck(:cell_name)
            cells = cells1 + cells2

            drugs1 = Drug.where("drug_name LIKE ?", "#{input}%").pluck(:drug_name)
            drugs2 = SourceDrugName.where("drug_name LIKE ?", "#{input}%").pluck(:drug_name)
            drugs = drugs1 + drugs2

            targets = Target.where("target_name LIKE ?", "#{input}%").pluck(:target_name)


            tissues1 = Tissue.where("tissue_name LIKE ?", "#{input}%").pluck(:tissue_name)
            tissues2 = SourceTissueName.where("tissue_name LIKE ?", "#{input}%").pluck(:tissue_name)
            tissues = tissues1 + tissues2

            datasets = Dataset.where("dataset_name LIKE ?", "#{input}%").pluck(:dataset_name)

            total_results = cells.uniq.count + drugs.uniq.count + targets.uniq.count + tissues.uniq.count + datasets.uniq.count 

            other_result = []

            if total_results == 0 || allowed_searches == "none"
                other_result << "No results to show!"
            end
            render json: other_result
        end
    end

end
