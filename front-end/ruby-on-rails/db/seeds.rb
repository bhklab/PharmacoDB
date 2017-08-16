require 'csv'

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'cell.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = Cell.new
# 	t.cell_id = row['cell_id']
# 	t.accession_id = row['accession_id']
# 	t.cell_name = row['cell_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'tissue.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = Tissue.new
# 	t.tissue_id = row['tissue_id']
# 	t.tissue_name = row['tissue_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'drug.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = Drug.new
# 	t.drug_id = row['drug_id']
# 	t.drug_name = row['drug_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'dataset.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = Dataset.new
# 	t.dataset_id = row['dataset_id']
# 	t.dataset_name = row['dataset_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'source.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = Source.new
# 	t.source_id = row['source_id']
# 	t.dataset_id = row['dataset_id']
# 	t.source_name = row['source_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'source_cell_name.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = SourceCellName.new
# 	t.cell_id = row['cell_id']
# 	t.source_id = row['source_id']
# 	t.cell_name = row['cell_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'source_tissue_name.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = SourceTissueName.new
# 	t.tissue_id = row['tissue_id']
# 	t.source_id = row['source_id']
# 	t.tissue_name = row['tissue_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'source_drug_name.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = SourceDrugName.new
# 	t.drug_id = row['drug_id']
# 	t.source_id = row['source_id']
# 	t.drug_name = row['drug_name']
# 	t.save
# end

# csv_text = File.read(Rails.root.join('lib', 'seeds', 'cell_tissue.csv'))
# csv = CSV.parse(csv_text, :headers => true, :encoding => 'ISO-8859-1')
# csv.each do |row|
# 	t = CellTissue.new
# 	t.cell_id = row['cell_id']
# 	t.tissue_id = row['tissue_id']
# 	t.save
# end

# summary
puts "There are now #{Cell.count} rows in the cells table"
puts "There are now #{Tissue.count} rows in the tissues table"
puts "There are now #{Drug.count} rows in the drugs table"
puts "There are now #{Dataset.count} rows in the datasets table"
puts "There are now #{Source.count} rows in the sources table"
puts "There are now #{SourceCellName.count} rows in the source_cell_names table"
puts "There are now #{SourceTissueName.count} rows in the source_tissue_names table"
puts "There are now #{SourceDrugName.count} rows in the source_drug_names table"
puts "There are now #{CellTissue.count} rows in the cell_tissues table"
puts "There are now #{Experiment.count} rows in the experiments table"
puts "There are now #{DoseResponse.count} rows in the dose_responses table"
puts "There are now #{Cellosauru.count} rows in the cellosaurus table"